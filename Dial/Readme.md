# 网络编程
### 服务端
- 建立并绑定 Socket：首先服务端使用 socket() 函数建立网络套接字，然后使用 bind() 函数为套接字绑定指定的 IP 和端口；
- 监听请求：接下来，服务端使用 listen() 函数监听客户端对绑定 IP 和端口的请求；
- 接收连接：如果有请求过来，并通过三次握手成功建立连接，则使用 accept() 函数接收并处理该连接；
- 处理请求与发送响应：服务端通过 read() 函数从上述已建立连接读取客户端发送的请求数据，经过处理后再通过 write() 函数将响应数据发送给客户端。
### 客户端
- 建立 Socket：客户端同样使用 socket()函数建立网络套接字；
- 建立连接：然后调用 connect() 函数传入 IP 和端口号建立与指定服务端网络程序的连接；
- 发送请求与接收响应：连接建立成功后，客户端就可以通过 write() 函数向服务端发送数据，并使用 read() 函数从服务端接收响应。
### Go中的Dail
```
func Dial(network, address string) (Conn, error) {
    var d Dialer
    return d.Dial(network, address)
}
```
network表示要传入的网络协议，address表示要传入的ip地址或者域名
支持的协议类型
- tcp、tcp4、tcp6     host:port 或 [host]:port
- udp、udp4、udp6
- ip、ipv4、ipv6         ip protocol number/protocol name
- 代表 Unix 通信域下的一种内部 socket 协议    address： file system path
    - unix：以 SOCK_STREAM 为 socket 类型。
    - unixgram：以 SOCK_DGRAM 为 socket 类型。
    - unixpacket: 以 SOCK_SEQPACKET 为 socket 类型。

### 额外知识点
- GO程序运行在纯命令行的模式下，该干什么全靠运行参数
os.Args使程序获取运行它时给出的参数，它的类型是[]string
- os.Exit(code)退出程序  0 success; 1~125 error;defer函数不执行
- bytes.NewBuffer(buf []byte) *Buffer   本函数用于创建一个用于读取已存在数据的buffer；也用于指定用于写入的内部缓冲的大小；
buf作为初始内容创建并初始化一个Buffer，应为指定具体容量但长度为0的切片。