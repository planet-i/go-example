package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	fileUrl1 := "http://192.168.51.206:30805/api/s3download/?time=1711532575300&filename=%E5%9B%BE%20(80).jpg&URL=http://www.datrix206.com:7480/aa7135485078d0509863339468732701/2024032714_btgwyhegujayy_lv0.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=c9ca82da496e4d3aff7eacceb2ddb710%2F20240327%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240327T094255Z&X-Amz-Expires=10800&X-Amz-SignedHeaders=host&response-content-type=application%2Foctet-stream&X-Amz-Signature=10a573e0c1621bccd160e5a0dae3c0801cacdc361c6e70269bf2852c86c88a16"
	if err := DownFileSimple("test.png", fileUrl1); err != nil {
		panic(err)
	}
	strURL := `http://192.168.51.206:30805/api/s3download/?time=1711592373321&filename=%E8%8B%B1%E8%AF%AD%E5%9F%BA%E7%A1%80%E9%98%85%E8%AF%BB(1)%20%E6%9D%8E%E9%93%B6%E7%BE%8E%20238%202023-09-15%2009_45_00-11_30_00_20240327180413053027.mp4&URL=http://www.datrix206.com:7480/aa7135485078d0509863339468732701/2024032718_btgiwefzujayy_lv0.mp4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=c9ca82da496e4d3aff7eacceb2ddb710%2F20240328%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240328T021933Z&X-Amz-Expires=10800&X-Amz-SignedHeaders=host&response-content-type=application%2Foctet-stream&X-Amz-Signature=b34648545473fe0715cc2936d2a5fdbe4975006125ee9a50dfee94661d3739`
	if err := DownBigFile("长津湖.mp4", strURL); err != nil {
		panic(err)
	}
}

// DownFileSimple会将url下载到本地文件，它会在下载时写入，而不是将整个文件加载到内存中。
func DownFileSimple(filepath string, url string) error {
	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件
	_, err = io.Copy(file, resp.Body)
	return err
}

func DownBigFile(filepath string, url string) error {
	// 创建临时文件
	file, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		file.Close()
		return err
	}
	defer resp.Body.Close()

	// 写入文件
	counter := &WriteCounter{}
	if _, err := io.Copy(file, io.TeeReader(resp.Body, counter)); err != nil {
		file.Close()
		return err
	}
	file.Close()

	// 关闭临时文件后重命名文件
	if err := os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

// 打印进度
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %d complete", wc.Total)
}
