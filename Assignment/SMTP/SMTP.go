package main

import (
	"encoding/base64"
	"fmt"
	"net"
)

func sendAndRecv(conn net.Conn, msg, vrfy string) {
	if msg != "" {
		conn.Write([]byte(msg))
	}
	recv := make([]byte, 1024)
	if vrfy == "" {
		return
	}
	_, err := conn.Read(recv)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	recvStr := string(recv)
	fmt.Println(recvStr)
	if recvStr[:3] != vrfy {
		fmt.Printf("error: %s reply not received from server.\n", vrfy)
		return
	}
}

func main() {
	// Choose a mail server (e.g. Google mail server) and call it mailserver
	mailserver := "smtp.qq.com:587"

	// Create socket called clientSocket and establish a TCP connection with mailserver
	clientSocket, err := net.Dial("tcp", mailserver)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer clientSocket.Close()

	sendAndRecv(clientSocket, "", "220")

	// Send HELO command and print server response.
	heloCommand := "HELO Alice\r\n"
	sendAndRecv(clientSocket, heloCommand, "250")

	// Send AUTH LOGIN command and print server response.
	authCommand := "AUTH LOGIN\r\n"
	sendAndRecv(clientSocket, authCommand, "334")
	mail := "1873978303@qq.com"
	// password := "swloksfjdvxufdch"
	password := "harkepzllstvfbgi"
	mailBase64 := base64.StdEncoding.EncodeToString([]byte(mail)) + "\r\n"
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password)) + "\r\n"
	sendAndRecv(clientSocket, mailBase64, "334")
	sendAndRecv(clientSocket, passwordBase64, "235")

	// Send MAIL FROM command and print server response.
	mailCommand := fmt.Sprintf("MAIL FROM:<%s>\r\n", mail)
	sendAndRecv(clientSocket, mailCommand, "250")

	// Send RCPT TO command and print server response.
	recpCommand := fmt.Sprintf("RCPT TO:<%s>\r\n", mail)
	sendAndRecv(clientSocket, recpCommand, "250")

	// Send DATA command and print server response.
	dataCommand := "DATA\r\n"
	sendAndRecv(clientSocket, dataCommand, "354")

	// Send message data.
	msg := "I love computer networks!"
	content := fmt.Sprintf("from:%s\r\nto:%s\r\nsubject:%s\r\n\r\n%s", mail, mail, "this is a demo", msg)
	sendAndRecv(clientSocket, content, "")

	// Message ends with a single period.
	endmsg := "\r\n.\r\n"
	sendAndRecv(clientSocket, endmsg, "250")

	// Send QUIT command and get server response.
	quitCommand := "QUIT\r\n"
	sendAndRecv(clientSocket, quitCommand, "221")
}
