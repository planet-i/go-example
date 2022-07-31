package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("connect to server error : ", err.Error())
	}
	defer conn.Close()
	for {
		buf := make([]byte, 2)
		buf[0] = 10
		buf[1] = 2
		conn.Write(buf)
		recv := make([]byte, 2048)
		n, err2 := conn.Read(recv)
		if err2 != nil {
			fmt.Println("err2", err2)
			return
		}
		if n == 4 {
			data := recv[0:n]

			v := binary.BigEndian.Uint32(data)
			fmt.Println("V uint32 = ", v)
		}
		time.Sleep(1 * time.Second)
	}
}
