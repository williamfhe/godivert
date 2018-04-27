package godivert

type Direction bool

const (
	PacketBufferSize   = 1500
	PacketChanCapacity = 256

	WinDivertDirectionOutbound Direction = false
	WinDivertDirectionInbound  Direction = true
)

func (d Direction) String() string {
	if bool(d) {
		return "Inbound"
	}
	return "Outbound"
}
