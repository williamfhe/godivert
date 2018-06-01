package godivert

import (
	"fmt"
	"github.com/williamfhe/godivert/header"
	"net"
)

// Represents a packet
type Packet struct {
	Raw       []byte
	Addr      *WinDivertAddress
	PacketLen uint

	IpHdr      header.IPHeader
	NextHeader header.ProtocolHeader

	ipVersion      int
	hdrLen         int
	nextHeaderType uint8

	parsed bool
}

// Parse the packet's headers
func (p *Packet) ParseHeaders() {
	p.ipVersion = int(p.Raw[0] >> 4)
	if p.ipVersion == 4 {
		p.hdrLen = int((p.Raw[0] & 0xf) << 2)
		p.nextHeaderType = p.Raw[9]
		p.IpHdr = header.NewIPv4Header(p.Raw)
	} else {
		p.hdrLen = 40
		p.nextHeaderType = p.Raw[6]
		p.IpHdr = header.NewIPv6Header(p.Raw)
	}

	switch p.nextHeaderType {
	case header.ICMPv4:
		p.NextHeader = header.NewICMPv4Header(p.Raw[p.hdrLen : p.hdrLen+header.ICMPv4HeaderLen])
	case header.TCP:
		p.NextHeader = header.NewTCPHeader(p.Raw[p.hdrLen:])
	case header.UDP:
		p.NextHeader = header.NewUDPHeader(p.Raw[p.hdrLen : p.hdrLen+header.UDPHeaderLen])
	case header.ICMPv6:
		p.NextHeader = header.NewICMPv6Header(p.Raw[p.hdrLen : p.hdrLen+header.ICMPv6HeaderLen])
	default:
		// Protocol not implemented
		p.NextHeader = nil
	}

	p.parsed = true
}

func (p *Packet) String() string {
	p.VerifyParsed()

	nextHeaderType := p.NextHeaderType()
	return fmt.Sprintf("Packet {\n"+
		"\tIPHeader=%v\n"+
		"\tNextHeaderType=(%d)->%s\n"+
		"\tNextHeader: %v\n"+
		"\tWinDivertAddr=%v\n"+
		"\tRawData=%v\n"+
		"}",
		p.IpHdr, nextHeaderType, header.ProtocolName(nextHeaderType), p.NextHeader, p.Addr, p.Raw)
}

// Returns the version of the IP protocol
// Shortcut for ipHdr.Version()
func (p *Packet) IpVersion() int {
	return p.ipVersion
}

// Returns the IP Protocol number of the next Header
// https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
func (p *Packet) NextHeaderType() uint8 {
	p.VerifyParsed()

	return p.nextHeaderType
}

// Returns the source IP of the packet
// Shortcut for IpHdr.SrcIP()
func (p *Packet) SrcIP() net.IP {
	p.VerifyParsed()

	return p.IpHdr.SrcIP()
}

// Sets the source IP of the packet
// Shortcut for IpHdr.SetSrcIP()
func (p *Packet) SetSrcIP(ip net.IP) {
	p.VerifyParsed()

	p.IpHdr.SetSrcIP(ip)
}

// Returns the destination IP of the packet
// Shortcut for IpHdr.DstIP()
func (p *Packet) DstIP() net.IP {
	p.VerifyParsed()

	return p.IpHdr.DstIP()
}

// Sets the destination IP of the packet
// Shortcut for IpHdr.SetDstIP()
func (p *Packet) SetDstIP(ip net.IP) {
	p.VerifyParsed()

	p.IpHdr.SetDstIP(ip)
}

// Returns the source port of the packet
// Shortcut for NextHeader.SrcPort()
func (p *Packet) SrcPort() (uint16, error) {
	p.VerifyParsed()

	if p.NextHeader == nil {
		return 0, fmt.Errorf("cannot get source port on protocolID=%d, protocol not implemented", p.nextHeaderType)
	}

	return p.NextHeader.SrcPort()
}

// Sets the source port of the packet
// Shortcut for NextHeader.SetSrcPort()
func (p *Packet) SetSrcPort(port uint16) error {
	p.VerifyParsed()

	if p.NextHeader == nil {
		return fmt.Errorf("cannot change source port on protocolID=%d, protocol not implemented", p.nextHeaderType)
	}

	return p.NextHeader.SetSrcPort(port)

}

// Returns the destination port of the packet
// Shortcut for NextHeader.DstPort()
func (p *Packet) DstPort() (uint16, error) {
	p.VerifyParsed()

	if p.NextHeader == nil {
		return 0, fmt.Errorf("cannot change get port on protocolID=%d, protocol not implemented", p.nextHeaderType)
	}

	return p.NextHeader.DstPort()
}

// Sets the destination port of the packet
// Shortcut for NextHeader.SetDstPort()
func (p *Packet) SetDstPort(port uint16) error {
	p.VerifyParsed()

	if p.NextHeader == nil {
		return fmt.Errorf("cannot change destination port on protocolID=%d, protocol not implemented", p.nextHeaderType)
	}

	return p.NextHeader.SetDstPort(port)
}

// Returns the name of the protocol
func (p *Packet) NextHeaderProtocolName() string {
	return header.ProtocolName(p.NextHeaderType())
}

// Inject the packet on the Network Stack
// If the packet has been modified calls WinDivertHelperCalcChecksum to get a new checksum
func (p *Packet) Send(wd *WinDivertHandle) (uint, error) {
	if p.parsed && (p.IpHdr.NeedNewChecksum() || p.NextHeader != nil && p.NextHeader.NeedNewChecksum()) {
		wd.HelperCalcChecksum(p)
	}
	return wd.Send(p)
}

// Recalculate the packet's checksum
// Shortcut for WinDivertHelperCalcChecksum
func (p *Packet) CalcNewChecksum(wd *WinDivertHandle) {
	wd.HelperCalcChecksum(p)
}

// Check if the headers have already been parsed and call ParseHeaders() if not
func (p *Packet) VerifyParsed() {
	if !p.parsed {
		p.ParseHeaders()
	}
}

// Returns the Direction of the packet
// WinDivertDirectionInbound (true) for inbound Packets
// WinDivertDirectionOutbound (false) for outbound packets
// Shortcut for Addr.Direction()
func (p *Packet) Direction() Direction {
	return p.Addr.Direction()
}

// Check the packet with the filter
// Returns true if the packet matches the filter
func (p *Packet) EvalFilter(filter string) (bool, error) {
	return HelperEvalFilter(p, filter)
}
