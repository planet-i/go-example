//网络协议开发中经常需要将int、float、bool等转为二进制数据，float32、64 与[]byte处理：

//1、写出uint64类型的变量a和字节数组互转的代码
//2、写出int64和string互转的代码
package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"
)

func main() {
	var a uint64
	a = 1
	b := make([]byte, 8)
	d := make([]byte, 8)

	//float转int或uint
	f := 1.5
	g := int(f)
	gg := *(*uint32)(unsafe.Pointer(&f))
	fmt.Println(f, g, gg)

	// uint64转[]byte
	// 用标准库带的函数
	binary.BigEndian.PutUint64(b, a)
	binary.LittleEndian.PutUint64(d, a)
	fmt.Println(a, b, d)
	// 位操作 大端序
	b[0] = byte(a >> 56)
	b[1] = byte(a >> 48)
	b[2] = byte(a >> 40)
	b[3] = byte(a >> 32)
	b[4] = byte(a >> 24)
	b[5] = byte(a >> 16)
	b[6] = byte(a >> 8)
	b[7] = byte(a)
	// []byte转uint64
	//用标准库的函数
	a = binary.BigEndian.Uint64(b)
	c := binary.LittleEndian.Uint64(d)
	fmt.Println(a, c, b, d)
	// 位操作
	a = uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	c = uint64(d[0]) | uint64(d[1])<<8 | uint64(d[2])<<16 | uint64(d[3])<<24 |
		uint64(d[4])<<32 | uint64(d[5])<<40 | uint64(d[6])<<48 | uint64(d[7])<<56
	fmt.Println(a, c, b, d)
	//取int的字节
	bb := 256
	aa := byte(bb)
	fmt.Println(aa)
	// string与[]byte转换
	ss := "abc"
	bbb := []byte(ss)
	sss := string(bbb)
	fmt.Println(bbb, sss)
	// string    int
	fmt.Println("string互转int")
	s := "12fv个机会6"
	i, _ := strconv.Atoi(s)
	i64, _ := strconv.ParseInt(s, 10, 64)
	ui64, _ := strconv.ParseUint(s, 10, 64)
	fmt.Println(s, i, i64, ui64)

	i = 123456
	s1 := strconv.Itoa(i)
	s2 := strconv.FormatInt(i64, 10)
	s3 := strconv.FormatUint(ui64, 10)
	fmt.Println(i, i64, s1, s2, s3)
	//string和int64转只能适用于字符串里全部是数字，
}

//IEEE 754 在类型中的应用
