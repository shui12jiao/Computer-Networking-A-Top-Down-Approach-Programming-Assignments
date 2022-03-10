package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
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

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("data error: ", err)
	}

	reg := regexp.MustCompile(`^([A-Z]+)\b /(\w+.\w+)?`)
	matches := reg.FindAllStringSubmatch(string(buf), 1)
	method := matches[0][1]
	fileName := matches[0][2]

	if method != "GET" {
		writeHeader(conn, 405, 0)
	}
	if fileName == "" {
		goto NotFound
	}

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

NotFound:
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
