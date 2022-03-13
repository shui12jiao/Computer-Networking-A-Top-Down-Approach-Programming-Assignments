package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func pingUDP(conn *net.UDPConn, ch chan time.Duration) {
	defer conn.Close()

	msg := "nyaru saikyou!"

	start := time.Now()
	conn.Write([]byte(msg))
	buf := make([]byte, 16)
	conn.SetReadDeadline(start.Add(time.Second))
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err)
			ch <- time.Hour
			return
		}
		if n != 0 {
			dur := time.Since(start)
			fmt.Printf("recv: %s\n", string(buf))
			ch <- dur
			return
		}
	}

}

func main() {
	socket := "localhost:80"
	if arg := os.Args[1]; arg != "" {
		reg := regexp.MustCompile(`^(localhost|(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d))?:\d+$`)
		if reg.MatchString(arg) {
			if arg[0] == ':' {
				socket = "localhost" + arg
			} else {
				socket = arg
			}
		} else {
			fmt.Println("wrong socket")
			os.Exit(1)
		}
	}
	times := 10
	if os.Args[2] != "" {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("arg times error: ", err)
		} else {
			times = n
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", socket)
	if err != nil {
		fmt.Println("socket error: ", err)
	}

	ch := make(chan time.Duration, times)
	durations := make([]time.Duration, times)

	for i := 0; i < times; i++ {
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Printf("dial %d error: %s\n", i+1, err.Error())
			continue
		}
		time.Sleep(time.Duration(time.Second))
		go pingUDP(conn, ch)
	}

	for i := 0; i < times; i++ {
		time, ok := <-ch
		if !ok {
			fmt.Println("something wrong with channel")
		}
		durations[i] = time
	}
	close(ch)

	for idx, dur := range durations {
		if dur > time.Second {
			fmt.Printf("%d. Reply from %s: bytes=14 timeout\n", idx+1, socket)
		} else if dur < time.Millisecond {
			fmt.Printf("%d. Reply from %s: bytes=14 time<1ms\n", idx+1, socket)
		} else {
			fmt.Printf("%d. Reply from %s: bytes=14 time=%dms\n", idx+1, socket, dur/1e6)
		}
	}
}

//TODO 格式
// Pinging baidu.com [220.181.38.251] with 32 bytes of data:
// Reply from 220.181.38.251: bytes=32 time=29ms TTL=53
// Reply from 220.181.38.251: bytes=32 time=29ms TTL=53
// Reply from 220.181.38.251: bytes=32 time=29ms TTL=53
// Reply from 220.181.38.251: bytes=32 time=31ms TTL=53

// Ping statistics for 220.181.38.251:
//     Packets: Sent = 4, Received = 4, Lost = 0 (0% loss),
// Approximate round trip times in milli-seconds:
//     Minimum = 29ms, Maximum = 31ms, Average = 29ms
