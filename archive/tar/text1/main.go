package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// tar包实现tar格式压缩文件的存取
func main() {
	// Create and add some files to the archive.
	var buf bytes.Buffer //bytes包的Buffer类型 代替io.Reader / io.Writer 类型

	// 创建一个写入 buf 的*Writer。
	tw := tar.NewWriter(&buf)
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling license."},
	}
	for _, file := range files {
		// Header代表tar档案文件里的单个头。
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)), // 为下一个文件写入多少字节
		}
		// 写入hdr并准备接受文件内容 若当前文件未完全写入，返回错误
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		// 向tar档案文件的当前记录中写入数据
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}
	}
	// Close关闭tar档案文件，会将缓冲中未写入下层的io.Writer接口的数据刷新到下层。
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	// Open and iterate through the files in the archive.

	// 创建一个从buf读取的Reader
	tr := tar.NewReader(&buf)
	for {
		//转入tar档案文件下一记录，它会返回下一记录的头域。
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}

}
