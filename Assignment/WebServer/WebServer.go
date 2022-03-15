package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

var fileNames []string = getFileNames("./file")

func getFileNames(path string) (names []string) {
	names = make([]string, 3)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("server error:", err)
	}

	for _, fi := range files {
		if fi.IsDir() {
			tmp := getFileNames(path + "/" + fi.Name())
			names = append(names, tmp...)
		}
		names = append(names, path+"/"+fi.Name())
	}
	return
}

func handleConn(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("data error:", err)
	}

	reg := regexp.MustCompile(`^([A-Z]+)\b (/.*) \b`)
	matches := reg.FindAllStringSubmatch(string(buf), 1)
	method := matches[0][1]
	fileName := matches[0][2]

	cType := "text/html"
	if method != "GET" {
		writeHeader(conn, 405, 0, cType)
		return
	}

	for _, name := range fileNames {
		if "."+fileName == name {
			file, err := ioutil.ReadFile(name)
			//TODO  增加file buffer
			//TODO  content-type分析

			if err != nil {
				fmt.Println("file error:", err)
			}
			len := len(file)

			tmp := strings.Split(fileName, ".")
			if tmp[1] != "" {
				cType = tmp[1]
			}
			writeHeader(conn, 200, len, cType)
			conn.Write(file)
			return
		}
	}

	writeHeader(conn, 404, 0, cType)
}

func writeHeader(conn *net.TCPConn, status int, length int, cType string) {
	buf := "HTTP/1.1 %d\r\nContent-Type:%s\r\nContent-Length:%d\r\n\r\n"
	header := []byte(fmt.Sprintf(buf, status, cType, length))
	conn.Write(header)
}

func main() {
	tcpServer, err := net.ResolveTCPAddr("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Println("port error:", err)
	}

	listen, err := net.ListenTCP("tcp", tcpServer)
	if err != nil {
		fmt.Println("listen error:", err)
	}
	defer listen.Close()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept error:", err)
		}
		go handleConn(conn)

	}
}
