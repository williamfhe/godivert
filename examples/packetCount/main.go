package main

import (
	"fmt"
	"time"
	"github.com/williamfhe/godivert"
	"github.com/williamfhe/godivert/header"
)

var icmpv4, icmpv6, udp, tcp, unknown, served uint

func checkPacket(wd *godivert.WinDivertHandle, packetChan  <- chan *godivert.Packet) {
	for packet := range packetChan {
		countPacket(packet)
		wd.Send(packet)
	}
}

func countPacket(packet *godivert.Packet) {
	served++
	switch packet.NextHeaderType() {
	case header.ICMPv4:
		icmpv4++
	case header.ICMPv6:
		icmpv6++
	case header.TCP:
		tcp++
	case header.UDP:
		udp++
	default:
		unknown++
	}
}


func main() {
	winDivert, err := godivert.NewWinDivertHandle("true")
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting")

	packetChan, err := winDivert.Packets()
	if err != nil {
		panic(err)
	}


	n := 50
	for i := 0; i < n; i++ {
		go checkPacket(winDivert, packetChan)

	}

	time.Sleep(15 * time.Second)

	winDivert.Close()
	fmt.Println("Stopping...")

	fmt.Printf("Served: %d packets\n", served)

	fmt.Printf("ICMPv4=%d ICMPv6=%d UDP=%d TCP=%d Unknown=%d", icmpv4, icmpv6, udp, tcp, unknown)
}
