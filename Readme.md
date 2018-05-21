# GoDivert

Go bindings for [WinDivert](https://github.com/basil00/Divert).

WinDivert is a user-mode packet capture-and-divert package for Windows.

## Installation

```bash
go get github.com/williamfhe/godivert
```

## Introduction

The binding's documentation can be found [Here](https://godoc.org/github.com/williamfhe/godivert).

If you don't have the **WinDivert dll** installed on your System or you want to load a specific **WinDivert dll** you should do :

```go
godivert.LoadDLL("PathToThe64bitDLL", "PathToThe32bitDLL")
```

The path can be a **relative path** to the *.exe* **current directory** or an **absolute path**.

Note that the driver must be in the **same directory** as the **dll**.
**LoadDLL** will then load the **dll** depending on your **OS architecture**.

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
packet.SetDstPort(1234) // Sets the destination port
packet.Send(winDivert) // Sends the packet back on the network stack
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

## Examples

### Capturing and Printing a Packet

```go
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

## Blocking Protocol by IP

```go
package main

import (
    "net"
    "time"
    "github.com/williamfhe/godivert"
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
    winDivert, err := godivert.NewWinDivertHandle("icmp")
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

Forbid all ICMP packets to reach 1.1.1.1 for 1 minute.

Try it :

```bash
ping 1.1.1.1
```

### Packet Count

```go
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

```

Count all protocols passing by for 15 seconds.