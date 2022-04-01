package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const ICMP_ECHO_REQUEST = 8

// def checksum(string):
// csum = 0
// countTo = (len(string) // 2) * 2
// count = 0
// while count < countTo:
// thisVal = ord(string[count+1]) * 256 + ord(string[count])
// csum = csum + thisVal
// csum = csum & 0xffffffff
// count = count + 2
// if countTo < len(string):
// csum = csum + ord(string[len(string) - 1])
// csum = csum & 0xffffffff
// csum = (csum >> 16) + (csum & 0xffff)
// csum = csum + (csum >> 16)
// answer = ~csum
// answer = answer & 0xffff
// answer = answer >> 8 | (answer << 8 & 0xff00)
// return answer

// func receiveOnePing(mySocket, ID, timeout, destAddr){
// 	timeLeft = timeout
// 	while 1:
// 	startedSelect = time.time()
// 	whatReady = select.select([mySocket], [], [], timeLeft)
// howLongInSelect = (time.time() - startedSelect)
// if whatReady[0] == []: # Timeout
// return "Request timed out."
// timeReceived = time.time()
// recPacket, addr = mySocket.recvfrom(1024)

//  #Fill in start

//  #Fetch the ICMP header from the IP packet

//  #Fill in end
// timeLeft = timeLeft - howLongInSelect
// if timeLeft <= 0:
// return "Request timed out."
// }

// func sendOnePing(mySocket, destAddr, ID){
// # Header is type (8), code (8), checksum (16), id (16), sequence (16)
// myChecksum = 0
// # Make a dummy header with a 0 checksum
// # struct -- Interpret strings as packed binary data
// header = struct.pack("bbHHh", ICMP_ECHO_REQUEST, 0, myChecksum, ID, 1)
// data = struct.pack("d", time.time())
// # Calculate the checksum on the data and the dummy header.
// myChecksum = checksum(str(header + data))
// # Get the right checksum, and put in the header
// if sys.platform == 'darwin':
// # Convert 16-bit integers from host to network byte order
// myChecksum = htons(myChecksum) & 0xffff
// else:
// myChecksum = htons(myChecksum)
// header = struct.pack("bbHHh", ICMP_ECHO_REQUEST, 0, myChecksum, ID, 1)
// packet = header + data
// mySocket.sendto(packet, (destAddr, 1)) # AF_INET address must be tuple, not str
// # Both LISTS and TUPLES consist of a number of objects
// # which can be referenced by their position number within the object.
// }

// func doOnePing(destAddr, timeout){
// icmp = getprotobyname("icmp")
// # SOCK_RAW is a powerful socket type. For more details: http://sockraw.org/papers/sock_raw
// mySocket = socket(AF_INET, SOCK_RAW, icmp)
// myID = os.getpid() & 0xFFFF # Return the current process i
// sendOnePing(mySocket, destAddr, myID)
// delay = receiveOnePing(mySocket, myID, timeout, destAddr)
// mySocket.close()
// return delay
// }

// func ping(host, timeout=1){
// # timeout=1 means: If one second goes by without a reply from the server,
// # the client assumes that either the client's ping or the server's pong is lost
// dest = gethostbyname(host)
// print("Pinging " + dest + " using Python:")
// print("")
// # Send ping requests to a server separated by approximately one second
// while 1 :
// delay = doOnePing(dest, timeout)
// print(delay)
// time.sleep(1)# one second
// return delay
// }

// ping("google.com")

type ICMP struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Sequence uint16
}

func checksum(raw []byte) uint16 {
	cks := uint32(0)
	len := 8
	idx := 0
	for idx < len {
		cks += uint32(raw[idx])<<8 + uint32(raw[idx+1])
		idx += 2
	}
	cks += (cks >> 16)
	return uint16(^cks)
}

func send(conn *net.IPConn) {
	hdr := ICMP{
		Type:     8,
		Code:     0,
		Checksum: 0,
		ID:       0,
		Sequence: 0,
	}
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, hdr)
	hdr.Checksum = checksum(buf.Bytes())

	conn.Write(buf.Bytes())
}

func receive(conn *net.IPConn) {
	buf := make([]byte, 256)
	conn.Read(buf)

}

func ping(conn *net.IPConn, timeout time.Duration) {
	defer conn.Close()
	send(conn)
	receive(conn)
}

func getIP() (laddr, raddr *net.IPAddr) {
	// address := "localhost"
	// if len(os.Args) > 1 {
	// 	address = os.Args[1]
	// }
	raddr, err := net.ResolveIPAddr("ip", "192.168.31.1") //TODO
	if err != nil {
		fmt.Printf("get raddr error:%s", err)
	}

	laddr, err = net.ResolveIPAddr("ip", "localhost") //TODO
	if err != nil {
		fmt.Printf("get laddr error:%s", err)
	}

	return
}

func main() {
	var (
		laddr, raddr = getIP()
		count        = 5
		timeout      = time.Second
	)
	conn, err := net.DialIP("ip:icmp", laddr, raddr)
	if err != nil {
		fmt.Printf("dial ip error:%s", err)
	}

	for i := 0; i < count; i++ {
		go ping(conn, timeout)
	}
}
