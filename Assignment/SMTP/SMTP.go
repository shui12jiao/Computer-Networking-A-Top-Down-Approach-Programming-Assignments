package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	msg := "\r\n I love computer networks!"
	endmsg := "\r\n.\r\n"

	// Choose a mail server (e.g. Google mail server) and call it mailserver
	mailserver := "smtp.qq.com:465"

	// Create socket called clientSocket and establish a TCP connection with mailserver
	clientSocket, err := smtp.Dial(mailserver)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer clientSocket.Close()
	clientSocket.Mail(msg)

	// recv := clientSocket.recv(1024).decode()
	// print(recv)
	// if recv[:3] != '220':
	//     print('220 reply not received from server.')

	// // Send HELO command and print server response.
	// heloCommand = 'HELO Alice\r\n'
	// clientSocket.send(heloCommand.encode())
	// recv1 = clientSocket.recv(1024).decode()
	// print(recv1)
	// if recv1[:3] != '250':
	//     print('250 reply not received from server.')

	// Send MAIL FROM command and print server response.
	// Fill in start
	// Fill in end

	// Send RCPT TO command and print server response.
	// Fill in start
	// Fill in end

	// Send DATA command and print server response.
	// Fill in start
	// Fill in end

	// Send message data.
	// Fill in start
	// Fill in end

	// Message ends with a single period.
	// Fill in start
	// Fill in end

	// Send QUIT command and get server response.
	// Fill in start
	// Fill in end
}
