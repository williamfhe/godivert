package header

import (
	"encoding/binary"
	"fmt"
)

// Represents a TCP header
// https://en.wikipedia.org/wiki/Transmission_Control_Protocol#TCP_segment_structure
type TCPHeader struct {
	Raw      []byte
	Modified bool
}

func NewTCPHeader(raw []byte) *TCPHeader {
	hdrLen := (raw[12] >> 4) * 4
	return &TCPHeader{
		Raw: raw[:hdrLen],
	}
}

func (h *TCPHeader) String() string {
	if h == nil {
		return "<nil>"
	}

	srcPort, _ := h.SrcPort()
	dstPort, _ := h.DstPort()

	return fmt.Sprintf("{\n"+
		"\t\tProtocol=TCP\n"+
		"\t\tSrcPort=%d\n"+
		"\t\tDstPort=%d\n"+
		"\t\tSeqNum=%#x\n"+
		"\t\tAckNum=%d\n"+
		"\t\tHeaderLen=%d\n"+
		"\t\tReserved=%d\n"+
		"\t\tFlags={NS=%t CWR=%t ECE=%t URG=%t ACK=%t PSH=%t RST=%t SYN=%t FIN=%t}\n"+
		"\t\tWindow=%d\n"+
		"\t\tCheckSum=%#x\n"+
		"\t\tUrgPtr=%d\n"+
		"\t\tOptions=%v\n"+
		"\t}\n",
		srcPort, dstPort, h.SeqNum(), h.AckNum(), h.HeaderLen(), h.Reserved(), h.NS(), h.CWR(), h.ECE(), h.URG(), h.ACK(), h.PSH(), h.RST(), h.SYN(), h.FIN(), h.Window(), h.Checksum(), h.UrgPtr(), h.Options())
}

// Reads the header's bytes and returns the source port
func (h *TCPHeader) SrcPort() (uint16, error) {
	return binary.BigEndian.Uint16(h.Raw[0:2]), nil
}

// Reads the header's bytes and returns the destination port
func (h *TCPHeader) DstPort() (uint16, error) {
	return binary.BigEndian.Uint16(h.Raw[2:4]), nil
}

// Sets the source port
func (h *TCPHeader) SetSrcPort(port uint16) error {
	h.Modified = true
	h.Raw[0] = uint8(port >> 8)
	h.Raw[1] = uint8(port & 0xff)
	return nil
}

// Sets the destination port
func (h *TCPHeader) SetDstPort(port uint16) error {
	h.Modified = true
	h.Raw[2] = uint8(port >> 8)
	h.Raw[3] = uint8(port & 0xff)
	return nil
}

// Reads the header's bytes and returns the sequence number
func (h *TCPHeader) SeqNum() uint32 {
	return binary.BigEndian.Uint32(h.Raw[4:8])
}

// Reads the header's bytes and returns the acknowledgment number
func (h *TCPHeader) AckNum() uint32 {
	return binary.BigEndian.Uint32(h.Raw[8:12])
}

// Reads the header's bytes and returns the length of the header in bytes
func (h *TCPHeader) HeaderLen() int {
	return int(h.DataOffset()) * 4
}

// Reads the header's bytes and returns the data offset
func (h *TCPHeader) DataOffset() uint8 {
	return h.Raw[12] >> 4
}

// Reads the header's bytes and returns the reserved part
func (h *TCPHeader) Reserved() uint8 {
	return (h.Raw[12] >> 1) & 0x7
}

// FLAGS START

// Reads the header's bytes and returns the NS flag as a boolean
func (h *TCPHeader) NS() bool {
	return h.Raw[12]&0x1 == 1
}

// Reads the header's bytes and returns the CWR as a boolean
func (h *TCPHeader) CWR() bool {
	return h.Raw[13]>>7 == 1
}

// Reads the header's bytes and returns the ECE flag as a boolean
func (h *TCPHeader) ECE() bool {
	return (h.Raw[13]>>6)&0x1 == 1
}

// Reads the header's bytes and returns the URG flag as a boolean
func (h *TCPHeader) URG() bool {
	return (h.Raw[13]>>5)&0x1 == 1
}

// Reads the header's bytes and returns the ACK flag as a boolean
func (h *TCPHeader) ACK() bool {
	return (h.Raw[13]>>4)&0x1 == 1
}

// Reads the header's bytes and returns the PSH flag as a boolean
func (h *TCPHeader) PSH() bool {
	return (h.Raw[13]>>3)&0x1 == 1
}

// Reads the header's bytes and returns the RST flag as a boolean
func (h *TCPHeader) RST() bool {
	return (h.Raw[13]>>2)&0x1 == 1

}

// Reads the header's bytes and returns the SYN flag as a boolean
func (h *TCPHeader) SYN() bool {
	return (h.Raw[13]>>1)&0x1 == 1
}

// Reads the header's bytes and returns the FIN flag as a boolean
func (h *TCPHeader) FIN() bool {
	return h.Raw[13]&0x1 == 1
}

// END FLAGS

// Reads the header's bytes and returns the window size
func (h *TCPHeader) Window() uint16 {
	return binary.BigEndian.Uint16(h.Raw[14:16])
}

// Reads the header's bytes and returns the checksum
func (h *TCPHeader) Checksum() uint16 {
	return binary.BigEndian.Uint16(h.Raw[16:18])
}

// Reads the header's bytes and returns the urgent pointer field
func (h *TCPHeader) UrgPtr() uint16 {
	return binary.BigEndian.Uint16(h.Raw[18:20])
}

// Reads the header's bytes and returns the options as a byte slice if they exist or nil
func (h *TCPHeader) Options() []byte {
	hdrLen := h.HeaderLen()
	if hdrLen <= TCPHeaderLen {
		return nil
	}
	return h.Raw[TCPHeaderLen:hdrLen]
}

// Returns true if the header has been modified
func (h *TCPHeader) NeedNewChecksum() bool {
	return h.Modified
}
