package godivert

import "fmt"

// WinDivertAddress : Represents a WinDivertAddress struct
// See : https://reqrypt.org/windivert-doc.html#divert_address
// As go doesn't not support bit fields
// we use a little trick to get the Direction, Loopback, Import and PseudoChecksum fields
type WinDivertAddress struct {
	Timestamp int64
	IfIdx     uint32
	SubIfIdx  uint32
	Data      uint8
}

func (w *WinDivertAddress) String() string {
	return fmt.Sprintf("{\n"+
		"\t\tTimestamp=%d\n"+
		"\t\tInteface={IfIdx=%d SubIfIdx=%d}\n"+
		"\t\tDirection=%v\n"+
		"\t\tLoopback=%t\n"+
		"\t\tImpostor=%t\n"+
		"\t\tPseudoChecksum={IP=%t TCP=%t UDP=%t}\n"+
		"\t}",
		w.Timestamp, w.IfIdx, w.SubIfIdx, w.Direction(), w.Loopback(), w.Impostor(),
		w.PseudoIPChecksum(), w.PseudoTCPChecksum(), w.PseudoUDPChecksum())
}

// Direction : Returns the direction of the packet
// WinDivertDirectionInbound (true) for inbounds packets
// WinDivertDirectionOutbounds (false) for outbounds packets
func (w *WinDivertAddress) Direction() Direction {
	return Direction(w.Data&0x1 == 1)
}

// Loopback : Returns true if the packet is a loopback packet
func (w *WinDivertAddress) Loopback() bool {
	return (w.Data>>1)&0x1 == 1
}

// Impostor : Returns true if the packet is an impostor
// See https://reqrypt.org/windivert-doc.html#divert_address for more information
func (w *WinDivertAddress) Impostor() bool {
	return (w.Data>>2)&0x1 == 1
}

// PseudoIPChecksum : Returns true if the packet uses a pseudo IP checksum
func (w *WinDivertAddress) PseudoIPChecksum() bool {
	return (w.Data>>3)&0x1 == 1
}

// PseudoIPChecksum : Returns true if the packet uses a pseudo TCP checksum
func (w *WinDivertAddress) PseudoTCPChecksum() bool {
	return (w.Data>>4)&0x1 == 1
}

// PseudoIPChecksum : Returns true if the packet uses a pseudo UDP checksum
func (w *WinDivertAddress) PseudoUDPChecksum() bool {
	return (w.Data>>5)&0x1 == 1
}
