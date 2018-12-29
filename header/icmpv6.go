package header

import (
	"fmt"
)

import (
	"encoding/binary"
	"errors"
)

// Represents a ICMPv6 header
// https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol_for_IPv6#Message_types_and_formats
type ICMPv6Header struct {
	Raw      []byte
	Modified bool
}

func NewICMPv6Header(raw []byte) *ICMPv6Header {
	return &ICMPv6Header{
		Raw: raw,
	}
}

func (h *ICMPv6Header) String() string {
	if h == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{\n"+
		"\t\tProtocol=ICMPv6\n"+
		"\t\tType=%d\n"+
		"\t\tCode=%d\n"+
		"\t\tHeaderLen=%d\n"+
		"\t\tChecksum=%#x\n"+
		"\t\tBody=%#x\n"+
		"\t}", h.Type(), h.Code(), ICMPv6HeaderLen, h.Checksum(), h.Body())
}

// Reads the header's bytes and returns the type
func (h *ICMPv6Header) Type() uint8 {
	return h.Raw[0]
}

func (h *ICMPv6Header) SetType(hType uint8) {
	h.Modified = true
	h.Raw[0] = hType
}

// Reads the header's bytes and returns the type
func (h *ICMPv6Header) Code() uint8 {
	return h.Raw[1]
}

func (h *ICMPv6Header) SetCode(code uint8) {
	h.Modified = true
	h.Raw[1] = code
}

// Reads the header's bytes and returns the body (Rest of the Header)
func (h *ICMPv6Header) Body() uint32 {
	return binary.BigEndian.Uint32(h.Raw[4:8])
}

func (h *ICMPv6Header) SetBody(body uint32) {
	h.Modified = true
	h.Raw[4] = uint8(body >> 24)
	h.Raw[5] = uint8((body >> 16) & 0xff)
	h.Raw[6] = uint8((body >> 8) & 0xff)
	h.Raw[7] = uint8(body & 0xff)
}

// Reads the header's bytes and returns the Checksum
func (h *ICMPv6Header) Checksum() uint16 {
	return binary.BigEndian.Uint16(h.Raw[2:4])
}

// Returns the length of the header in bytes (8 bytes)
func (h *ICMPv6Header) HeaderLen() int {
	return ICMPv6HeaderLen
}

// Not used for ICMPv6
func (h *ICMPv6Header) SrcPort() (uint16, error) {
	return 0, errors.New("ICMPv6 header has no source port")
}

// Not used for ICMPv6
func (h *ICMPv6Header) DstPort() (uint16, error) {
	return 0, errors.New("ICMPv6 header has no destination port")
}

// Not used for ICMPv6
func (h *ICMPv6Header) SetSrcPort(port uint16) error {
	return errors.New("ICMPv6 header has no source port")
}

// Not used for ICMPv6
func (h *ICMPv6Header) SetDstPort(port uint16) error {
	return errors.New("ICMPv6 header has no destination port")
}

// Returns true if the header has been modified
func (h *ICMPv6Header) NeedNewChecksum() bool {
	return h.Modified
}
