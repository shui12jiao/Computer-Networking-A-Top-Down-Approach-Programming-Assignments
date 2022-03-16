package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func dialServer(url string) []byte {
	res := strings.SplitN(url, "/", 2)
	address := res[0]

	// hasPort := false
	// for _, c := range address {
	// 	if c == ':' {
	// 		hasPort = true
	// 		break
	// 	}
	// }
	// if !hasPort {
	// 	address += ":80"
	// }
	// file := res[1]

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("dial server error:", err)
		return nil
	}
	defer conn.Close()

	writeReqHeader(conn, url)

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("server error:", err)
		return nil
	}
	return buf
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
	method := matches[0][1] //--------------------go on
	url := matches[0][2]
	cType := "text/html"
	if len(url) < 3 {
		writeRespHeader(conn, 404, 0, cType)
		return
	}
	if method != "GET" {
		writeRespHeader(conn, 405, 0, cType)
		return
	}

	file := dialServer(url)
	if file == nil {
		fmt.Println("get file error:", err)
		return
	}

	conn.Write(file)
	fmt.Println(string(file))
}

func writeRespHeader(conn net.Conn, status int, length int, cType string) {
	buf := "HTTP/1.1 %d\r\nContent-Type:%s\r\nContent-Length:%d\r\n\r\n"
	header := []byte(fmt.Sprintf(buf, status, cType, length))
	conn.Write(header)
}

func writeReqHeader(conn net.Conn, adr string) {
	header := []byte(fmt.Sprintf("GET http://%s HTTP/1.1\r\n", adr))
	conn.Write(header)
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

// GET /www.google.com HTTP/1.1
// Host: localhost:8080
// Connection: keep-alive
// Cache-Control: max-age=0
// sec-ch-ua: " Not A;Brand";v="99", "Chromium";v="99", "Microsoft Edge";v="99"
// sec-ch-ua-mobile: ?0
// sec-ch-ua-platform: "Windows"
// Upgrade-Insecure-Requests: 1
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 Edg/99.0.1150.39
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
// Sec-Fetch-Site: none
// Sec-Fetch-Mode: navigate
// Sec-Fetch-User: ?1
// Sec-Fetch-Dest: document
// Accept-Encoding: gzip, deflate, br
// Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6
