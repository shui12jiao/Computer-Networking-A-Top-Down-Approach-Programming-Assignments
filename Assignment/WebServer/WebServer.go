package main

import (
	"fmt"
	"io/ioutil"
	"net"
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

func HandleConn(conn net.TCPConn)
	defer conn.Close()
	buf := conn.Read()
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

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Printf("accept error: ", err)
		}
		go HandleConn(conn)

	}
	fmt.Println(fileNames)
}
