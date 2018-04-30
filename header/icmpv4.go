package header

import (
	"fmt"
)

import (
	"encoding/binary"
	"errors"
)

// Represents a ICMP header
// https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol#Header
type ICMPv4Header struct {
	Raw      []byte
	Modified bool
}

func NewICMPv4Header(raw []byte) *ICMPv4Header {
	return &ICMPv4Header{
		Raw: raw,
	}
}

func (h *ICMPv4Header) String() string {
	if h == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{\n"+
		"\t\tProtocol=ICMPv4\n"+
		"\t\tType=%d\n"+
		"\t\tCode=%d\n"+
		"\t\tHeaderLen=%d\n"+
		"\t\tChecksum=%#x\n"+
		"\t\tBody=%#x\n"+
		"\t}", h.Type(), h.Code(), ICMPv4HeaderLen, h.Checksum(), h.Body())
}

// Reads the header's bytes and returns the type
func (h *ICMPv4Header) Type() uint8 {
	return h.Raw[0]
}

func (h *ICMPv4Header) SetType(hType uint8) {
	h.Modified = true
	h.Raw[0] = hType
}

// Reads the header's bytes and returns the code
func (h *ICMPv4Header) Code() uint8 {
	return h.Raw[1]
}

func (h *ICMPv4Header) SetCode(code uint8) {
	h.Modified = true
	h.Raw[1] = code
}

// Reads the header's bytes and returns the body (Rest of the Header)
func (h *ICMPv4Header) Body() uint32 {
	return binary.BigEndian.Uint32(h.Raw[4:8])
}

func (h *ICMPv4Header) SetBody(body uint32) {
	h.Modified = true
	h.Raw[4] = uint8(body >> 24)
	h.Raw[5] = uint8((body >> 16) & 0xff)
	h.Raw[6] = uint8((body >> 8) & 0xff)
	h.Raw[7] = uint8(body & 0xff)
}

// Reads the header's bytes and returns the Checksum
func (h *ICMPv4Header) Checksum() uint16 {
	return binary.BigEndian.Uint16(h.Raw[2:4])
}

// Returns the length of the header in bytes (8 bytes)
func (h *ICMPv4Header) HeaderLen() uint8 {
	return ICMPv4HeaderLen
}

// Not used for ICMP
func (h *ICMPv4Header) SrcPort() (uint16, error) {
	return 0, errors.New("ICMPv4 header has no source port")
}

// Not used for ICMP
func (h *ICMPv4Header) DstPort() (uint16, error) {
	return 0, errors.New("ICMPv4 header has no destination port")
}

// Not used for ICMP
func (h *ICMPv4Header) SetSrcPort(port uint16) error {
	return errors.New("ICMPv4 header has no source port")
}

// Not used for ICMP
func (h *ICMPv4Header) SetDstPort(port uint16) error {
	return errors.New("ICMPv4 header has no destination port")
}

// Returns true if the header has been modified
func (h *ICMPv4Header) NeedNewChecksum() bool {
	return h.Modified
}
