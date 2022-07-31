# net
net包对于网络I/O提供了便携式接口，包括TCP/IP,UDP，域名解析以及Unix Socket。尽管net包提供了大量访问底层的接口，但是大多数情况下，客户端仅仅只需要最基本的接口，例如Dial，LIsten，Accepte以及分配的conn连接和listener接口。 crypto/tls包使用相同的接口以及类似的Dial和Listen函数。

### http
http服务器：主要在于服务端（server）如何接受客户端（clinet）的 request，并向client返回response。
Clinet -> Requests ->  [Multiplexer(router) -> handler  -> Response -> Clinet
- Multiplexer
    - Golang中的Multiplexer基于ServeMux结构，同时也实现了Handler接口。
- Handle
    - handler函数： 具有func(w http.ResponseWriter, r *http.Requests)签名的函数
    - handler处理器函数: 经过HandlerFunc结构包装的handler函数，它是实现了ServeHTTP接口方法的函数。调用- handler处理器的ServeHTTP方法时，即调用handler函数本身。
    - handler对象：实现了Handler接口ServeHTTP方法的结构。
一个处理器就是一个拥有ServeHTTP方法的接口
