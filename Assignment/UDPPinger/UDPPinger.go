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
			// fmt.Println("read error:", err)
			ch <- time.Hour
			return
		}
		if n != 0 {
			dur := time.Since(start)
			// fmt.Printf("recv: %s\n", string(buf))
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
			fmt.Println("arg times error:", err)
		} else {
			times = n
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", socket)
	if err != nil {
		fmt.Println("socket error:", err)
	}

	ch := make(chan time.Duration, times)

	for i := 0; i < times; i++ {
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Printf("dial %d error:%s\n", i+1, err.Error())
			continue
		}
		go pingUDP(conn, ch)
	}

	var avgTime time.Duration
	var maxTime time.Duration
	var minTime time.Duration = time.Hour
	var recvNum int
	fmt.Printf("Pinging [%s] with 14 bytes of data:\n", socket)
	for i := 0; i < times; i++ {
		dur, ok := <-ch
		if !ok {
			fmt.Println("something wrong with channel")
		}

		if dur > time.Second {
			fmt.Println("Request timed out.")
		} else if dur < time.Millisecond {
			fmt.Printf("Reply from %s: bytes=14 time<1ms\n", socket)
			recvNum++
		} else {
			fmt.Printf("Reply from %s: bytes=14 time=%dms\n", socket, dur/1e6)
			recvNum++
			avgTime += dur
			if dur > maxTime {
				maxTime = dur
			}
			if dur < minTime {
				minTime = dur
			}
		}
	}
	close(ch)
	fmt.Printf("\nPing statistics for [%s]:\n", socket)
	fmt.Printf("    Packets: Sent = %d, Received = %d, Lost = %d (%d%% loss)\n", times, recvNum, times-recvNum, (times-recvNum)*100/times)
	fmt.Println("Approximate round trip times in milli-seconds:")
	fmt.Printf("    Minimum = %dms, Maximum = %dms, Average = %dms", minTime/1e6, maxTime/1e6, avgTime/time.Duration(recvNum)/1e6)
}
