package godivert

type Direction bool

const (
	PacketBufferSize   = 1500
	PacketChanCapacity = 256

	WinDivertDirectionOutbound Direction = false
	WinDivertDirectionInbound  Direction = true
)

const (
	WinDivertFlagSniff uint8 = 1 << iota
	WinDivertFlagDrop  uint8 = 1 << iota
	WinDivertFlagDebug uint8 = 1 << iota
)

func (d Direction) String() string {
	if bool(d) {
		return "Inbound"
	}
	return "Outbound"
}
