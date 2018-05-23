package main

import (
	"fmt"
	"github.com/williamfhe/godivert"
)

func main() {
	winDivert, err := godivert.NewWinDivertHandle("true")
	if err != nil {
		panic(err)
	}
	defer winDivert.Close()

	packet, err := winDivert.Recv()
	if err != nil {
		panic(err)
	}

	fmt.Println(packet)

    packet.Send(winDivert)
}
