package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/net/html/charset"

	//"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Record struct {
	Name string `csv:"source term"`
	Age  string `csv:"target term"`
}

func testCSVReaderGBK() {
	f, err := os.Open("test.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	var records []Record
	if err := gocsv.Unmarshal(f, &records); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(records)
}

// 测试CSV文件的读取
func testCSVReader() {
	// 替换为实际的 CSV 文件路径
	filePath := "test.csv"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("打开文件失败: %s", err)
	}
	defer file.Close()

	// 编码检测
	decoder := GetCsvTransformer(file)

	// 初始化Reader
	reader := csv.NewReader(transform.NewReader(file, decoder))
	//reader := csv.NewReader(bufio.NewReader(file))
	//reader := csv.NewReader(file)
	// 允许非标准引号
	reader.LazyQuotes = true
	reader.ReuseRecord = true

	// 读取 CSV 文件内容
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("读取 CSV 文件失败: %s", err)
	}

	// 输出读取的内容
	for _, record := range records {
		log.Println(strings.Join(record, ","))
	}
}

//	func readFirstFewBytes(file multipart.File) []byte {
//		const size = 512 // 检测编码可能不需要读取整个文件，512字节通常足够
//		buffer := make([]byte, size)
//		n, err := file.Read(buffer)
//		if err != nil && err != io.EOF { // 如果不是正常结束，则返回错误
//			return nil
//		}
//		// 跳过 BOM 头
//		if n >= 3 && bytes.Equal(buffer[:3], []byte{0xEF, 0xBB, 0xBF}) {
//			buffer = buffer[3:]
//			n -= 3
//		}
//		return buffer[:n]
//	}
func GetCsvTransformer(file multipart.File) transform.Transformer {
	// 编码检测
	encoding, name, _ := charset.DetermineEncoding(readFirstFewBytes(file), "")
	// 动态选择解码器
	var decoder transform.Transformer
	switch name {
	case "GBK", "GB18030", "windows-1252":
		decoder = simplifiedchinese.GBK.NewDecoder()
	default:
		log.Printf("未知编码 %s，尝试UTF-8解码", name)
		decoder = encoding.NewDecoder()
	}
	return decoder
}

func readFirstFewBytes(file multipart.File) []byte {
	const size = 1024
	buffer := make([]byte, size)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil
	}
	data := buffer[:n]

	// 检测并跳过所有常见 BOM 头
	bomOffsets := detectBOM(data)
	if bomOffsets > 0 {
		data = data[bomOffsets:]
	}

	// 关键：重置文件指针到 BOM 之后的位置，避免后续读取数据缺失
	_, _ = file.Seek(int64(bomOffsets), io.SeekStart)

	return data
}

// 检测并返回 BOM 长度（支持 UTF-8/16/32）
func detectBOM(data []byte) int {
	if len(data) >= 4 {
		switch {
		case bytes.HasPrefix(data, []byte{0x00, 0x00, 0xFE, 0xFF}):
			return 4 // UTF-32 BE
		case bytes.HasPrefix(data, []byte{0xFF, 0xFE, 0x00, 0x00}):
			return 4 // UTF-32 LE
		}
	}
	if len(data) >= 3 && bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		return 3 // UTF-8
	}
	if len(data) >= 2 {
		switch {
		case bytes.HasPrefix(data, []byte{0xFE, 0xFF}):
			return 2 // UTF-16 BE
		case bytes.HasPrefix(data, []byte{0xFF, 0xFE}):
			return 2 // UTF-16 LE
		}
	}
	return 0
}
