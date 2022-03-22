package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

func urlHandle(url string) (address, fileName string) {
	res := strings.SplitN(url, "/", 2)
	address = res[0]
	if !strings.ContainsRune(address, ':') {
		address += ":80"
	}
	if len(res) < 2 {
		fileName = "index.html"
	} else {
		fileName = res[1]
	}
	return
}

func dialServer(url string) []byte {
	address, fileName := urlHandle(url)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("dial server error:", err)
		return nil
	}
	defer conn.Close()

	url = address + fileName
	file := FileBuffer.searchAndUpdateBuffer(url)
	if file == nil {
		conn.Write(reqHeader(address, fileName))

		var buf bytes.Buffer
		_, err = io.Copy(&buf, conn)
		if err != nil {
			fmt.Println("server error:", err)
			return nil
		}
		err = FileBuffer.addFile(url, buf.Bytes())
		if err != nil {
			fmt.Println("add file to buffer error:", err)
		}
		return buf.Bytes()
	} else {
		return file
	}
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
	cType := "text/html"
	if len(matches) < 1 || len(matches[0]) < 2 {
		conn.Write(respHeader(501, 0, cType))
		return
	}
	method := matches[0][1]
	url := matches[0][2]
	if len(url) < 3 {
		conn.Write(respHeader(501, 0, cType))
		return
	}
	if method != "GET" {
		conn.Write(respHeader(405, 0, cType))
		return
	}

	response := dialServer(url)
	if response == nil {
		conn.Write(respHeader(404, 0, cType))
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

func reqHeader(host, fileName string) (header []byte) {
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
