package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ICMP struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Sequence uint16
}

func checksum(raw []byte) uint16 {
	cks := uint32(0)
	len := len(raw)
	idx := 0
	for len-idx > 1 {
		cks += uint32(raw[idx])<<8 + uint32(raw[idx+1])
		idx += 2
	}
	if len-idx == 1 {
		cks += uint32(raw[idx])
	}
	cks += (cks >> 16)
	return uint16(^cks)
}

func verify(raw []byte) bool {
	sum := uint32(0)
	len := len(raw)
	idx := 0
	for len-idx < 1 {
		sum += uint32(raw[idx])<<8 + uint32(raw[idx+1])
		idx += 2
	}
	if len-idx == 1 {
		sum += uint32(raw[idx])
	}

	return sum == uint32(0)
}

func unmarshal(raw []byte) (hdr ICMP, data string) {
	hdr = ICMP{
		Type:     uint8(raw[0]),
		Code:     uint8(raw[1]),
		Checksum: uint16(raw[2])<<8 + uint16(raw[3]),
		ID:       uint16(raw[4])<<8 + uint16(raw[5]),
		Sequence: uint16(raw[6])<<8 + uint16(raw[7]),
	}
	data = string(raw[8:])
	return
}

func send(conn *net.IPConn, seq int) {
	hdr := ICMP{
		Type:     8,
		Code:     0,
		Checksum: 0,
		ID:       1,
		Sequence: uint16(seq),
	}
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, hdr)
	hdr.Checksum = checksum(buf.Bytes())
	buf.Reset()
	binary.Write(buf, binary.BigEndian, hdr)
	conn.Write(buf.Bytes())
}

func receive(conn *net.IPConn, seq int) (delay time.Duration, ttl int) {
	startTime := time.Now()
	done := make(chan time.Duration, 1)
	buf := make([]byte, 64)
	go func() {
		conn.Read(buf)
		done <- time.Since(startTime)
	}()

	select {
	case delay = <-done:
		buf = func(b []byte) []byte {
			if len(b) < 20 {
				return b
			}
			len := (b[0] & 0x0f) << 2
			ttl = int(b[8])
			b = b[len:]
			return b
		}(buf)

		right := verify(buf)
		if !right {
			fmt.Println("checksum wrong")
			return
		}
		hdr, _ := unmarshal(buf)
		if hdr.Sequence != uint16(seq) || hdr.Code != 0 || hdr.Type != 0 {
			fmt.Println("wrong reply")
		}
	case <-time.After(time.Second):
		fmt.Println("timeout")
		delay = time.Second
	}
	return
}

func ping(conn *net.IPConn, timeout time.Duration, seq int, recp *[]int) {
	send(conn, seq)
	delay, ttl := receive(conn, seq)
	if delay >= time.Second {
		fmt.Println("Request timed out.")
	} else if delay >= time.Millisecond {
		delay /= 1e6
		fmt.Printf("Reply from %s: time=%dms TTL=%d\n", conn.RemoteAddr().String(), delay, ttl)
		(*recp)[0] += 1
		*recp = append(*recp, int(delay))
	} else {
		fmt.Printf("Reply from %s: time<1ms TTL=%d\n", conn.RemoteAddr().String(), ttl)
		(*recp)[0] += 1
	}
}

func getIP() (laddr, raddr *net.IPAddr) {
	// address := "localhost"
	// if len(os.Args) > 1 {
	// 	address = os.Args[1]
	// }
	raddr, err := net.ResolveIPAddr("ip", "www.baidu.com") //TODO
	if err != nil {
		fmt.Println("get raddr error:", err)
	}

	laddr, err = net.ResolveIPAddr("ip", "218.194.53.147") //TODO
	if err != nil {
		fmt.Println("get laddr error:", err)
	}

	return
}

func main() {
	var (
		laddr, raddr = getIP()
		count        = 3
		timeout      = time.Second
	)
	conn, err := net.DialIP("ip:icmp", laddr, raddr)
	if err != nil {
		fmt.Println("dial ip error:", err)
		return
	}
	defer conn.Close()

	recp := make([]int, 1, count+1)
	fmt.Printf("Pinging baidu.com [%s]:\n", conn.RemoteAddr().String())
	for i := 0; i < count; i++ {
		ping(conn, timeout, i, &recp)
		// go ping(conn, timeout)--使用多线程需要(加锁|sleep|通道)
	}
	min, max, avg := 1000, 0, 0
	for _, v := range recp[1:] {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
		avg += v
	}
	avg /= (len(recp) - 1)
	fmt.Printf("\nPing statistics for %s:\n", conn.RemoteAddr().String())
	fmt.Printf("    Packets: Sent = %d, Received = %d, Lost = %d (%d%% loss),\n", count, recp[0], count-recp[0], (count-recp[0])/count*100)
	fmt.Println("Approximate round trip times in milli-seconds:") //TODO
	fmt.Printf("    Minimum = %dms, Maximum = %dms, Average = %dms", min, max, avg)
}
