package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

func pingUDP(conn *net.UDPConn) {

}

func main() {
	socket := "localhost:80"
	if arg := os.Args[1]; arg == "" {
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
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("dial error: ", err)
	}

	for i := 0; i < times; i++ {
		go pingUDP(conn)
	}

}
