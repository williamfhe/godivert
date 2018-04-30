# GoDivert

Go bindings for [WinDivert](https://github.com/basil00/Divert).

## Installation

```bash
go get github.com/williamfhe/godivert
```

## Example

```go
package main

import (
    "github.com/williamfhe/godivert"
    "fmt"
)

func main() {
    winDivert, err := godivert.NewWinDivertHandle("true")
    if err != nil {
        panic(err)
    }

    packet, err := winDivert.Recv()
    if err != nil {
        panic(err)
    }

    fmt.Println(packet)

    packet.Send(winDivert)

    winDivert.Close()
}
```

Wait for a packet and print it.

```go
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
```

Forbid all packets to reach 1.1.1.1