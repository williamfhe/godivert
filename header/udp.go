package header

import (
	"encoding/binary"
	"fmt"
)

// UDPHeader : Represents a UDP header
// See : https://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
type UDPHeader struct {
	Raw      []byte
	Modified bool
}

func NewUDPHeader(raw []byte) *UDPHeader {
	return &UDPHeader{
		Raw: raw,
	}
}

func (h *UDPHeader) String() string {
	if h == nil {
		return "<nil>"
	}

	srcPort, _ := h.SrcPort()
	dstPort, _ := h.DstPort()

	return fmt.Sprintf("{\n"+
		"\t\tProtocol=UDP\n"+
		"\t\tSrcPort=%d\n"+
		"\t\tDstPort=%d\n"+
		"\t\tHeaderLen=%d\n"+
		"\t\tContentLen=%d\n"+
		"\t\tChecksum=%#x\n"+
		"\t}", srcPort, dstPort, h.HeaderLen(), h.Len(), h.Checksum())
}

// HeaderLen : Returns the length of the header in bytes (8 bytes)
func (h *UDPHeader) HeaderLen() uint8 {
	return UDPHeaderLen
}

// SrcPort : Reads the header's bytes and returns the source port
func (h *UDPHeader) SrcPort() (uint16, error) {
	return binary.BigEndian.Uint16(h.Raw[0:2]), nil
}

// DstPort : Reads the header's bytes and returns the destination port
func (h *UDPHeader) DstPort() (uint16, error) {
	return binary.BigEndian.Uint16(h.Raw[2:4]), nil
}

// SetSrcPort : Sets the source port
func (h *UDPHeader) SetSrcPort(port uint16) error {
	h.Modified = true
	h.Raw[0] = uint8(port >> 8)
	h.Raw[1] = uint8(port & 0xff)
	return nil
}

// SetDstPort : Sets the destination port
func (h *UDPHeader) SetDstPort(port uint16) error {
	h.Modified = true
	h.Raw[2] = uint8(port >> 8)
	h.Raw[3] = uint8(port & 0xff)
	return nil
}

// Len : Reads the header and returns the length of UDP header and UDP data in bytes
func (h *UDPHeader) Len() uint16 {
	return binary.BigEndian.Uint16(h.Raw[4:6])
}

// Checksum : Reads the header's bytes and returns the checksum
func (h *UDPHeader) Checksum() uint16 {
	return binary.BigEndian.Uint16(h.Raw[6:8])
}

// NeedNewChecksum : Returns true if the header has been modified
func (h *UDPHeader) NeedNewChecksum() bool {
	return h.Modified
}
