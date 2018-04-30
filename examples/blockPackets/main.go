package main

import (
	"github.com/williamfhe/godivert"
	"net"
	"time"
)

var cloudflareDNS = net.ParseIP("1.1.1.1")

func checkPacket(wd *godivert.WinDivertHandle, packetChan <-chan *godivert.Packet) {
	for packet := range packetChan {
		if !packet.DstIP().Equal(cloudflareDNS) {
			packet.Send(wd)
		}
	}
}

func main() {
	winDivert, err := godivert.NewWinDivertHandle("true")
	if err != nil {
		panic(err)
	}

	packetChan, err := winDivert.Packets()
	if err != nil {
		panic(err)
	}

	go checkPacket(winDivert, packetChan)

	time.Sleep(1 * time.Minute)
	winDivert.Close()
}
