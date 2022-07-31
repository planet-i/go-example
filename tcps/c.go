package tcps

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"time"
)

type Client struct {
	con  net.Conn
	Name string
	Uuid string
	Id   int64
}

func NewClient(con net.Conn) *Client {
	client := new(Client)
	client.con = con
	client.Uuid = ""
	return client
}

func (c *Client) write() {
	for {
		runtime.Gosched()
	}
}

func (c *Client) Wrtie(buf []byte) bool {
	n, err := c.con.Write(buf)
	if err != nil {
		fmt.Println("write error ")
		return false
	}
	if n > 0 {
		return true
	}
	return false
}

func (c *Client) WriteString(str string) bool {
	return c.Wrtie([]byte(str))
}

func (c *Client) read() {
	for {
		c.con.SetReadDeadline(time.Now().Add(60 * 5 * time.Second))
		buff := make([]byte, 2048)
		count, err := c.con.Read(buff)
		if err != nil {
			break
		}
		if count > 0 {
			data := buff[0:count]
			// recvie data
			fmt.Println(data)
			// jiexi data shu
			//v := uint32(data[0]) * uint32(data[1])
			v := rand.Uint32()
			fmt.Println("V ", v)
			bb := make([]byte, 4)
			binary.BigEndian.PutUint32(bb, v)
			c.Wrtie(bb)
		}
		runtime.Gosched()
	}
}

func (c *Client) Run() {
	go c.read()
	go c.write()
}
