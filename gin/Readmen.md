# GIN的学习
## Test1
### Gin的用处
- Gin是一个非常优秀的Golang Web Framework，弥补了net/http的不足，同时增加了很多日常Web开发使用的功能。
- Gin允许我们对特定的HTTP方法进行不同的处理，为我们实现RESTful API提供了方便
### 使用Gin
- 基于GOPATH开发：先使用go get -u github.com/gin-gonic/gin 下载gin，然后import导入
- GO Module方式：使用import直接导入，然后在go run运行时会自动下载gin包编译使用；也可以通过go mod tidy 来下载依赖的模块

### HTTP Method
- GET：用于获取服务器上的资源，就是在浏览器中直接输入网址回车请求的方法。
- POST：给服务端提交一个资源，导致服务器的资源发生变化。
- PUT：用请求有效载荷替换目标资源的所有当前表示。
- DELETE：删除指定的资源。
- PATCH：用于对资源应用部分修改。
- HEAD：请求一个与GET请求的响应相同的响应，但没有响应体.
- CONNECT；建立一个到由目标资源标识的服务器的隧道。
- OPTIONS：用于描述目标资源的通信选项。
- TRACE：沿着到目标资源的路径执行一个消息环回测试。

### RESTful API 规范
开发Web应用服务或者程序，其实就是对服务器的资源的CRUD（创建、检索、更新和删除）  
RESTful API 的规范建议我们使用特定的HTTP方法来对服务器上的资源进行操作。
- GET，表示读取服务器上的资源
    - HTTP GET https://www.flysnow.org/users      获取所有用户的信息
    - HTTP GET https://www.flysnow.org/users/123  获取id为123用户的信息。
- POST，表示在服务器上创建资源
    - HTTP POST https://www.flysnow.org/users    创建一个用户，会通过POST给服务器提供创建这个用户所需的全部信息
- PUT,表示更新或者替换服务器上的资源，操作的是单个资源
    - HTTP PUT https://www.flysnow.org/users/123    更新/替换id为123的这个用户，在更新的时候，会通过PUT提供更新这个用户需要的全部用户信息。
- DELETE，表示删除服务器上的资源，操作的是单个资源
    - HTTP DELETE https://www.flysnow.org/users/123 删除id为123的这个用户
- PATCH，表示更新/修改资源的一部分 操作的是单个资源
    - HTTP PATCH https://www.flysnow.org/users/123 更新/替换id为123的这个用户，在更新的时候，会通过PATCH提供更新这个用户需要更新的部分用户信息。

### 知识点
- c.JSON()  
func (c *Context) JSON(code int, obj interface{})  
参数：code是返回的HTTP Status Code; obj是内容  
- gin.H其实是一个map[string]interface{},声明为H类型  
type H map[string]interface{}

## Test2 
### Gin的路由
Gin的路由采用的是httprouter，所以它的路由参数的定义和httprouter也是一样的。  
Gin的路由是单一的，不能有重复  
Gin路径中的匹配都是字符串，它是不区分数字、字母和汉字的，都匹配。  
gin.RedirectTrailingSlash 等于 true 时 若指定路径没匹配路由，会自动重定向到有匹配路由的路径，将其赋值为false则不会自动重定向
### 路由参数
：匹配当前路径     * 匹配所有    
- c.Param()  获取定义的路由参数的值

### 查询参数  
URL查询参数，或者也可以简称为URL参数，是存在于我们请求的URL中，以?为起点，后面的k=v&k1=v1&k2=v2这样的字符串就是查询参数；在URL中，多个查询参数键值对通过&相连。  
- 示例：https://www.flysnow.org/search?q=golang&sitesearch=https%3A%2F%2Fwww.flysnow.org
- 查询参数：?q=golang&sitesearch=https%3A%2F%2Fwww.flysnow.org  
- 两个查询参数键值对:  
    q=golang  
    sitesearch=https%3A%2F%2Fwww.flysnow.org
#### Gin获取查询参数
- c.Query()：获取参数key的value。key不存在，则返回""字符串  
- c.DefaultQuery():会判断对应的key是否存在，如果不存在的话，则返回默认defaultValue值。  
- c.QueryArray():适用于?a=b&a=c&a=d情况，返回参数a对应的value数组

## Form表单
### Form
| **查询参数**  | **Form表单**     | **说明**                                    |
| ------------- | ---------------- | --------------------------------------- |
| Query         | PostForm         | 获取key对应的值，不存在为空字符串       |
| GetQuery      | GetPostForm      | 多返回一个key是否存在的结果             |
| QueryArray    | PostFormArray    | 获取key对应的数组，不存在返回一个空数组 |
| GetQueryArray | GetPostFormArray | 多返回一个key是否存在的结果             |
| QueryMap      | PostFormMap      | 获取可以对应的map，不存在返回空的map    |
| GetQueryMap   | GetPostFormMap   | 多返回一个key是否存在的结果             |
| DefaultQuery  | DefaultPostForm  | key不存在的话，可以指定返回的默认值     |


