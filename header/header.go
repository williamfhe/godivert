package header

import "net"

// Represents a IPv4 or IPv6 Header
type IPHeader interface {
	String() string

	Version() int
	HeaderLen() uint8
	NextHeader() uint8
	SrcIP() net.IP
	DstIP() net.IP
	SetSrcIP(net.IP)
	SetDstIP(net.IP)
	Checksum() (uint16, error)
	NeedNewChecksum() bool
}

// Represents a protocol header
// Supported headers are TCP, UDP, ICMPv4, ICMPv6
type ProtocolHeader interface {
	String() string

	HeaderLen() int
	Checksum() uint16
	SrcPort() (uint16, error)
	DstPort() (uint16, error)
	SetSrcPort(uint16) error
	SetDstPort(uint16) error
	NeedNewChecksum() bool
}

// Returns the name of the given protocol number
// See : https://en.wikipedia.org/wiki/List_of_IP_protocol_numbers
func ProtocolName(protocol uint8) string {
	switch protocol {
	case ICMPv4:
		return "ICMPv4"
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	case ICMPv6:
		return "ICMPv6"
	default:
		return "Unimplemented Protocol"
	}
}
