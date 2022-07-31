### encoding/json
通过 encoding/json 标准库生成和解析 JSON 格式数据
#### func Marshal(v interface{}) ([]byte, error)
参数为空接口，意味着可以传入各种类型数据。编码成功则返回对应的JSON格式文本，否则，返回错误信息。
channel、complex、函数类型的值不能编码进json
- 布尔值转化为 JSON 后还是布尔类型；

- 浮点数和整型会被转化为 JSON数字类型；

- 字符串将以UTF-8编码转化输出为Unicode字符集的字符串，特殊字符将会被转义,"<"和">"会转义为"\u003c"和"\u003e","&"转义为"\u0026".

- 数组和切片会转化为 JSON数组，但 []byte 类型的值将会被转化为 Base64 编码后的字符串，slice的零值编码为 null；

- 结构体会转化为 JSON 对象，并且只有结构体里边以大写字母开头的可被导出的字段才会被转化输出，而这些可导出的字段会作为 JSON 对象的字符串索引；

- 转化一个 map 类型的数据结构时，该数据的类型必须是 map[string]T（T 可以是 encoding/json 包支持的任意数据类型）。

#### func Unmarshal(data []byte, v interface{}) error
第一个参数为待解码的JSON格式文本，第二个参数为解码结果映射的Go 语言目标类型数据结构
在 JSON 结构已知情况下的解码
- 如果 JSON 中的字段在 v 中不存在，json.Unmarshal() 函数在解码过程中会丢弃该字段.
这个特性让我们可以从同一段 JSON 数据中筛选指定的值填充到多个不同的 Go 语言类型中。
- 对于 JSON 中没有而 v 中定义的字段，会以对应数据类型的默认值填充，
Json中数据结构的类型转换
- 布尔值将会转换为 Go 语言的 bool 类型；
- 数值会被转换为 Go 语言的 float64 类型；
- 字符串转换后还是 string 类型；
- JSON 数组会转换为 []interface{} 类型；
- JSON 对象会转换为map[string]interface{} 类型；
- null 值会转换为 nil。

#### Json的流式读写
- Decoder
 - func NewDecoder(r io.Reader) *Decoder   创建一个从r读取并解码json对象的*Decoder
    - func (dec *Decoder) Decode(v interface{}) error  从输入流读取下一个json编码值并保存在v指向的值里
- Encoder
 - func NewEncoder(w io.Writer) *Encoder   创建一个将数据写入w的*Encoder。
    - func (enc *Encoder) Encode(v interface{}) error   将v的json编码写入输出流，并会写入一个换行符
