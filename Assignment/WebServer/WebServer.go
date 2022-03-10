package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var fileNames []string = getFileNames("./file")

func getFileNames(path string) (names []string) {
	names = make([]string, 3)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("server error: ", err)
	}

	for _, fi := range files {
		if fi.IsDir() {
			tmp := getFileNames(path + "/" + fi.Name())
			names = append(names, tmp...)
		}
		names = append(names, fi.Name())
	}
	return
}

func handleConn(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 0, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("data error: ", err)
	}

	i := 0
	var method_bt strings.Builder
	for i < n && buf[i] != ' ' {
		method_bt.WriteByte(buf[i])
		i++
	}
	method := method_bt.String()
	if method != "GET" {
		writeHeader(conn, 405, 0)
		return
	}

	for i < n && buf[i] == ' ' {
		i++
	}
	var url_bt strings.Builder
	for i < n && buf[i] != ' ' {
		url_bt.WriteByte(buf[i])
		i++
	}
	url := strings.Split(url_bt.String(), "/")
	fileName := url[1]
	for _, name := range fileNames {
		if name == fileName {
			file, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Println("file error: ", err)
			}
			len := len(file)
			writeHeader(conn, 200, len)
			conn.Write(file)
			return
		}
	}
	writeHeader(conn, 404, 0)

}

func writeHeader(conn *net.TCPConn, status int, length int) {
	buf := "HTTP/1.1 %d\r\nContent-Type:text/html\r\nContent-Length:%d\r\n\r\n"
	header := []byte(fmt.Sprintf(buf, status, length))

	conn.Write(header)
}

func main() {
	tcpServer, err := net.ResolveTCPAddr("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Println("port error: ", err)
	}

	listen, err := net.ListenTCP("tcp", tcpServer)
	if err != nil {
		fmt.Println("listen error: ", err)
	}
	defer listen.Close()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
		}
		go handleConn(conn)

	}
}
