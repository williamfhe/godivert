# GoDivert

Go bindings for [WinDivert](https://github.com/basil00/Divert).

WinDivert is a user-mode packet capture-and-divert package for Windows.

## Installation

```bash
go get github.com/williamfhe/godivert
```

## Introduction

The binding's documentation can be found [Here](https://godoc.org/github.com/williamfhe/godivert).

To start create a new instance of **WinDivertHandle** by calling **NewWinDivertHandle** and passing the filter as a parameter.

Documentation of the **filter** can be found [Here](https://reqrypt.org/windivert-doc.html#filter_language).

```go
winDivert, err := godivert.NewWinDivertHandle("Your filter here")
```

**WinDivertHandle** is struct that you can use to call WinDivert's function like **Recv** or **Send**.

You can divert a packet from the network stack by using **winDivert.Recv()** where **winDivert** is an instance of **WinDivertHandle**.

```go
packet, err := winDivert.Recv()
```

You can then choose to send the packet or modify it.

```go
packet.SetDstPort(1234) // Set the destination port
packet.Send(winDivert) // Send the packet back on the network stack
```

You can get and set values from the packet's header by using the **_header_** package. Documentation on this package can be found [Here](https://godoc.org/github.com/williamfhe/godivert/header)
.

As the packet has been modified the **checksums** have to be recalculated before sending it back on the network stack.

It is done automatically if the packet has been modified when calling **packet.Send** but you can do it manually by calling **packet.CalcNewChecksum**.

To receive packets you can also use **winDivert.Packets**.

```go
packetChan, err := winDivert.Packets()
```

Here **_packetChan_** is a channel of **\*godivert.Packet** coming directly from the network stack.

Note that all packets diverted are guaranteed to match the filter given in **godivert.NewWinDivertHandle("You filter here")**

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