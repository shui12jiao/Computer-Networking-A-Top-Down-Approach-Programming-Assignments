package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func dialServer(url string) []byte {
	res := strings.SplitN(url, "/", 2)
	if len(res) < 2 {
		return nil
	}
	address := res[0]
	fileName := res[1]

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("dial server error:", err)
		return nil
	}
	defer conn.Close()

	conn.Write(reqHeader(fileName, address))

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Println("server error:", err)
		return nil
	}

	return response
}

// var buffers map[string][]byte = make(map[string][]byte)

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error:", err)
	}
	reg := regexp.MustCompile(`^([A-Z]+)\b /([^\s]*) \b`)
	matches := reg.FindAllStringSubmatch(string(buf), 1)
	if len(matches[0]) < 2 {
		fmt.Println("error parameter")
	}
	method := matches[0][1] //--------------------go on
	url := matches[0][2]
	cType := "text/html"
	if len(url) < 3 {
		conn.Write(respHeader(404, 0, cType))
		return
	}
	if method != "GET" {
		conn.Write(respHeader(405, 0, cType))
		return
	}

	response := dialServer(url)
	if response == nil {
		fmt.Println("error parameter")
		return
	}

	conn.Write(response)
	// fmt.Println(string(file))
}

func respHeader(status int, length int, cType string) (header []byte) {
	str := "HTTP/1.0 %d\r\nContent-Type:%s\r\nContent-Length:%d\r\n\r\n"
	header = []byte(fmt.Sprintf(str, status, cType, length))
	return
}

func reqHeader(fileName, host string) (header []byte) {
	str := "GET /%s HTTP/1.0\r\nHost: %s\r\n"
	header = []byte(fmt.Sprintf(str, fileName, host))
	return
}

func main() {
	serverAddress := "localhost:8080"
	listen, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("proxy server error:", err)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
		}
		go handleConn(conn)
	}
}

// HTTP/1.0 200\r\nContent-Type:text/html\r\nContent-Length:84\r\n\r\n#include <iostream>\r\n\r\nint main() {\r\n    std::cout << \"Hello World\" << std::endl;\r\n}
// Parse Error: Expected HTTP/
