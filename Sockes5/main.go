package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	tcpConnId uint64
	udpConnId uint64
)

type proxy struct {
	conn *net.TCPConn
}

type uServ struct {
	prefix        string
	sigQuit       chan bool
	clientUdpAddr *net.UDPAddr
	dstMap        map[string]string
	sync.RWMutex
}

func main() {
	// Listen Tcp, support ipv4 only. @20171204
	service := "0.0.0.0:8083" //如果又可能是ipv6呢
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	// 输出错误信息，然后退出
	errCheckWithFatal("ResolveTcpAddr err:%v\n", err)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	errCheckWithFatal("ListenTCP err:%v\n", err)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		errCheckWithPrint("Warn >>> AcceptTCP err:%v\n", err)
		if err != nil {
			continue
		}

		proxy := NewProxy(tcpConn)
		go proxy.run()

	}
}

func NewProxy(conn *net.TCPConn) *proxy {
	return &proxy{
		conn: conn,
	}
}

func (p *proxy) run() {
	defer p.conn.Close()
	remoteAddrStr := p.conn.RemoteAddr().String()

	// read
	buf := make([]byte, 1024)
	n, err := p.conn.Read(buf[0:])
	errCheckWithPrint("Error >>> p.conn.Read err:%v\n", err)
	if err != nil {
		return
	}

	// get data bytes
	b := buf[0:n]

	//check Methods
	if !p.isMethodsOk(b) {
		log.Printf("WARN >>> from %s request methods invalid.\n", remoteAddrStr)
		return
	}

	n, err = p.conn.Read(buf[0:])
	errCheckWithPrint("Error >>> p.conn.Read err:%v\n", err)
	if err != nil {
		return
	}
	b = buf[0:n]
	// check username/password
	if !p.isAuthOk(b) {
		log.Printf("WARN >>> from %s auth fail.\n", remoteAddrStr)
		return
	}

	// resolve proxy CMD(tcp or udp)
	/*
	   +----+-----+-------+------+----------+--------+
	   |VER | CMD | RSV | ATYP | DST.ADDR | DST.PORT |
	   +----+-----+-------+------+----------+--------+
	   |  1 |  1  |X’00’|  1   | Variable |     2    |
	   +----+-----+-------+------+----------+--------+
	*/
	n, err = p.conn.Read(buf[0:])
	errCheckWithPrint("Error >>> p.conn.Read err:%v\n", err)
	if err != nil {
		return
	}
	b = buf[0:n]

	if !p.isCmdSupport(b[0:2]) {
		return
	}

	if !p.isAtypSupport(b[3:4]) {
		return
	}

	dstAddr, dstPort := p.getAddrPort(b[3:])
	if dstAddr == "" {
		log.Printf("WARN >>> get NULL dstAddr,return.\n")
		return
	}

	switch b[1] {
	case 0x01: //CONNECT -> TCP
		p.tcpProxy(dstAddr, dstPort)
		return
	case 0x02:
		log.Println("WARN >>> get BIND command, not support.")
	case 0x03: //UDP
		p.udpProxy(dstAddr, dstPort)
		return
	default:
		return
	}
}

func (p *proxy) isMethodsOk(b []byte) bool {
	// username/password only.
	/*
	   +----+----------+----------+
	   |VER | NMETHODS |  METHODS |
	   +----+----------+----------+
	   |  1 |     1    | 1 to 255 |
	   +----+----------+----------+
	*/
	if b[0] != 0x05 || (b[1] == 0x01 && b[2] != 0x02) {
		p.conn.Write([]byte{0x05, 0xff})
		return false
	}
	p.conn.Write([]byte{0x05, 0x02})
	return true
}

/* X’00’ NO AUTHENTICATION REQUIRED（不需要验证）
X’01’ GSSAPI
X’02’ USERNAME/PASSWORD（用户名密码）
X’03’ to X’7F’ IANA ASSIGNED
X’80’ to X’FE’ RESERVED FOR PRIVATE METHODS
X’FF’ NO ACCEPTABLE METHODS（都不支持，没法连接了） */

func (p *proxy) isAuthOk(b []byte) bool {
	//0x01 | 用户名长度（1字节）| 用户名（长度根据用户名长度域指定） | 口令长度（1字节） | 口令（长度由口令长度域指定）

	// b[0] some software use 0x01,but base protocal must use 0x05,so keeps it here.
	b0 := b[0]

	nameLens := int(b[1])
	name := string(b[2 : 2+nameLens])

	passLens := int(b[2+nameLens])
	pass := string(b[2+nameLens+1 : 2+nameLens+1+passLens])

	//auth lens
	/*
	   if nameLens != 32 || passLens != 48 {
	       p.conn.Write([]byte{b0, 0xff})
	       return false
	   }
	*/

	// auth count
	// test auth here.
	if name != "abc" || pass != "123" {
		p.conn.Write([]byte{b0, 0xff})
		return false
	}
	p.conn.Write([]byte{b0, 0x00})
	return true
}

func (p *proxy) isCmdSupport(b []byte) bool {
	if b[0] != 0x05 {
		p.conn.Write([]byte{0x05, 0x01})
		return false
	}

	switch b[1] {
	case 0x01: //CONNECT
		return true
	case 0x02: //BIND
		//not support here.
		p.conn.Write([]byte{0x05, 0x07})
		return false
	case 0x03: //UDP
		return true
	default:
		p.conn.Write([]byte{0x05, 0x07})
		return false
	}
}

func (p *proxy) isAtypSupport(b []byte) bool {

	switch b[0] {
	case 0x01: //ipv4
		return true
	case 0x03: //domain
		return true
	case 0x04: //ipv6
		p.conn.Write([]byte{0x05, 0x08})
		return false
	default:
		p.conn.Write([]byte{0x05, 0x08})
		return false
	}
}

func (p *proxy) getAddrPort(b []byte) (string, int) {
	switch b[0] {
	case 0x01: //ipv4
		addr := net.IPv4(b[1], b[2], b[3], b[4]).String()
		port := int(b[5])*256 + int(b[6])
		return addr, port

	case 0x03: //domain
		domainLens := int(b[1])
		domain := string(b[2 : 2+domainLens])
		port := int(b[2+domainLens])*256 + int(b[2+domainLens+1])
		return domain, port
	default:
		return "", 0
	}
}

func (p *proxy) tcpProxy(dstAddr string, dstPort int) {
	addr := fmt.Sprintf("%s:%d", dstAddr, dstPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	errCheckWithPrint("Error >>> net.ResolveTcpAddr err:%v\n", err)
	if err != nil {
		p.conn.Write([]byte{0x05, 0x03})
		return
	}

	tcpConn, err := TcpDial(nil, tcpAddr, 60*time.Second)
	format := fmt.Sprintf("Error >>> net.DialTCP %s err", addr)
	errCheckWithPrint(format+":%v\n", err)
	if err != nil {
		p.conn.Write([]byte{0x05, 0x03})
		return
	}
	defer tcpConn.Close()

	time.Sleep(10 * 1e6)
	p.conn.SetNoDelay(true)
	tcpConn.SetNoDelay(true)

	p.conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	tcpConnId++
	prefix := fmt.Sprintf("> tcp: %08x", tcpConnId)
	closeSig := make(chan bool, 0)
	go pipe("%s >>> %d bytes send\n", prefix, p.conn, tcpConn, closeSig)
	go pipe("%s <<< %d bytes recieve\n", prefix, tcpConn, p.conn, closeSig)
	<-closeSig
	return
}

func (p *proxy) udpProxy(clientAddr string, clientPort int) {
	//log.Printf("get client bind -> %s:%d\n", clientAddr, clientPort)

	// get localIp from p.conn
	localExitAddr, err := net.ResolveTCPAddr("tcp", p.conn.LocalAddr().String())
	errCheckWithFatal("ResolveTCPAddr err:%v\n", err)
	localExitIP := localExitAddr.IP.String()

	service := fmt.Sprintf("%s:%d", localExitIP, 0)
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	errCheckWithFatal("ResolveUDPAddr err:%v\n", err)

	udpListener, err := net.ListenUDP("udp", udpAddr)
	errCheckWithFatal("ListenUDP err:%v\n", err)
	defer udpListener.Close()
	//udpListener.SetPKTINFO()

	bindUdpAddr, err := net.ResolveUDPAddr("udp", udpListener.LocalAddr().String())
	errCheckWithFatal("ResolveUDPAddr err:%v\n", err)
	bindPort := bindUdpAddr.Port

	//版本 | 代理的应答 |　保留1字节　| 地址类型 | 代理服务器地址 | 绑定的代理端口
	bindMsg := []byte{0x05, 0x00, 0x00, 0x01}
	buffer := bytes.NewBuffer(bindMsg)
	binary.Write(buffer, binary.BigEndian, localExitAddr.IP.To4())
	binary.Write(buffer, binary.BigEndian, uint16(bindPort))
	//log.Printf("Local udp BIND >>> %s:%d\n", localExitIP, bindPort)
	//log.Printf("Send udp BIND Msg >>> %v\n", buffer.Bytes())
	p.conn.Write(buffer.Bytes())

	finish := make(chan bool, 0)
	go udpKeepAlive(p.conn, finish)
	go udpServ(udpListener, finish)
	<-finish
	return
}

func errCheckWithFatal(format string, err error) {
	if err != nil {
		log.Fatalf("Fatal >>>"+format, err)
	}
	return
}

func errCheckWithPrint(format string, err error) {
	if err != nil {
		log.Printf(format, err)
	}
	return
}

func TcpDial(localAddr, remoteAddr *net.TCPAddr, timeout time.Duration) (*net.TCPConn, error) {
	returned := false
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	type rst struct {
		tcn *net.TCPConn
		error
	}

	rst_chan := make(chan *rst, 0)
	go func() {
		tcpConn, err := net.DialTCP("tcp", localAddr, remoteAddr)
		if err != nil {
			goto Finish
		} else if returned {
			tcpConn.Close()
			return
		}
	Finish:
		rst_chan <- &rst{tcn: tcpConn, error: err}
	}()

	select {
	case <-ticker.C:
		returned = true
		return nil, errors.New("connect timeout")
	case result := <-rst_chan:
		if result.error != nil {
			return nil, result.error
		}
		return result.tcn, nil
	}
}

func pipe(format, prefix string, src, dst *net.TCPConn, closeSig chan bool) {
	buf := make([]byte, 0xff)
	for {
		n, err := src.Read(buf[0:])
		if err != nil {
			closeSig <- true
			return
		}
		b := buf[0:n]
		_, err = dst.Write(b)
		if err != nil {
			closeSig <- true
			return
		}
		//log.Printf(format, prefix, n)
	}
}

func udpKeepAlive(tcpConn *net.TCPConn, finish chan<- bool) {
	tcpConn.SetKeepAlive(true)
	buf := make([]byte, 1024)
	for {
		_, err := tcpConn.Read(buf[0:])
		if err != nil {
			finish <- true
			return
		}
	}
}

func udpServ(udpConn *net.UDPConn, finish chan<- bool) {
	udpConnId++

	udpServ := &uServ{
		prefix:  fmt.Sprintf("udp >>> #%08x", udpConnId),
		sigQuit: make(chan bool, 0),
		dstMap:  make(map[string]string),
	}

	go udpServ.read(udpConn)
	udpServ.monitor()
	finish <- true
	return
}

func (uS *uServ) monitor() {
	for {
		select {
		case <-uS.sigQuit:
			return
		}
	}
}

func (uS *uServ) read(udpConn *net.UDPConn) {
	buf := make([]byte, 2048)
	buf2 := make([]byte, 2048)
	for {
		n, udpAddr, err := udpConn.ReadFromUDP(buf[0:])
		if err != nil {
			//log.Printf("%s ReadFromUDP err:%v\n", uS.prefix, err)
			uS.sigQuit <- true
			return
		}

		if uS.clientUdpAddr == nil {
			uS.clientUdpAddr = udpAddr
		}

		b := buf[0:n]
		if udpAddr.IP.String() == uS.clientUdpAddr.IP.String() { // from client
			/*
			   +----+------+------+----------+----------+----------+
			   |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
			   +----+------+------+----------+----------+----------+
			   |  2 |   1  |   1  | Variable |     2    | Variable |
			   +----+------+------+----------+----------+----------+
			*/
			if b[2] != 0x00 {
				log.Printf("%s WARN: FRAG do not support.\n", uS.prefix)
				continue
			}

			switch b[3] {
			case 0x01: //ipv4
				dstAddr := &net.UDPAddr{
					IP:   net.IPv4(b[4], b[5], b[6], b[7]),
					Port: int(b[8])*256 + int(b[9]),
				}

				uS.Lock()
				if _, exist := uS.dstMap[dstAddr.String()]; !exist {
					uS.dstMap[dstAddr.String()] = string(b[0:10])
				}
				uS.Unlock()

				udpConn.WriteToUDP(b[10:], dstAddr)
				//log.Printf("%s b-> %v\n", uS.prefix, b)
				//log.Printf("%s data-> %v\n", uS.prefix, b[10:])
			case 0x03: //domain
				domainLens := int(b[4])
				domain := string(b[5 : 5+domainLens])
				ipAddr, err := net.ResolveIPAddr("ip", domain)
				if err != nil {
					log.Printf("%s Error -> domain %s dns query err:%v\n", uS.prefix, domain, err)
					continue
				}
				dstAddr := &net.UDPAddr{
					IP:   ipAddr.IP,
					Port: int(b[5+domainLens])*256 + int(b[6+domainLens]),
				}

				uS.Lock()
				if _, exist := uS.dstMap[dstAddr.String()]; !exist {
					uS.dstMap[dstAddr.String()] = string(b[0 : 7+domainLens])
				}
				uS.Unlock()

				udpConn.WriteToUDP(b[7+domainLens:], dstAddr)
				//log.Printf("%s b-> %v\n", uS.prefix, b)
				//log.Printf("%s data-> %v\n", uS.prefix, b[7+domainLens:])

			default:
				log.Printf("%s WARN: ATYP %v do not support.\n", uS.prefix, b[3])
				continue
			}
		} else { // from dst Server
			uS.RLock()
			if v, exist := uS.dstMap[udpAddr.String()]; exist {
				uS.RUnlock()
				head := []byte(v)
				headLens := len(head)
				copy(buf2[0:], head[0:headLens])
				copy(buf2[headLens:], b[0:])
				sendData := buf2[0 : headLens+n]
				udpConn.WriteToUDP(sendData, uS.clientUdpAddr)
				//log.Printf("%s <<< head-> %v\n", uS.prefix, head)
				//log.Printf("%s <<< b-> %v\n", uS.prefix, b)
				//log.Printf("%s <<< data-> %v\n", uS.prefix, sendData)
			} else {
				fmt.Printf("%s WARN -> %s not in dstMap.\n", uS.prefix, udpAddr.String())
				uS.RUnlock()
				continue
			}
		}
	}
	uS.sigQuit <- true
	return
}
