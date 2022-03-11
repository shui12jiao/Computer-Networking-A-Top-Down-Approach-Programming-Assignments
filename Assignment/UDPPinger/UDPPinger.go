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
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Print("read error: ", err)
		}

		dur := time.Since(start)
		if n != 0 || dur > time.Second {
			fmt.Printf("recv: %s\n", string(buf))
			ch <- dur
			return
		}
	}

}

func main() {
	socket := "localhost:80"
	if arg := os.Args[1]; arg != "" {
		reg := regexp.MustCompile(`^(localhost|(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.
		(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.
		(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.
		(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d))?:\d+$`)
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
		var durS string
		if dur > time.Second {
			durS = "timeout"
		} else {
			durS = strconv.Itoa(int(dur / 1e6))
		}
		fmt.Printf("%d. Reply from %s: bytes=14 time=%sms", idx+1, socket, durS)
	}
}
