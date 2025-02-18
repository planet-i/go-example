package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	fixRequireShipping()
	deleteDuplicateProduct()
}

func fixRequireShipping() {
	// 打开日志文件
	file, err := os.Open("BSPro.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建CSV文件
	csvFile, err := os.Create("BS.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	// 创建CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// 写入CSV表头
	writer.Write([]string{"序号", "商品店铺ID", "商品ID", "商品SPU"})

	// 正则表达式用于匹配数字部分
	re := regexp.MustCompile(`(\d+)-(\d+)\s+(\d+)-(\w+)`)

	// 逐行读取日志文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否包含错误信息
		if strings.Contains(line, "Product variant does not exist") {
			// 使用正则表达式提取数字部分
			matches := re.FindStringSubmatch(line)
			if len(matches) == 5 {
				// 写入CSV文件
				writer.Write([]string{matches[1], matches[2], matches[3], matches[4]})
			}
		}
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func deleteDuplicateProduct() error {
	// 打开输入文件
	input, err := os.Open("MK_AU.log")
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer input.Close()

	// 创建输出文件
	output, err := os.Create("MK_AU.csv")
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer output.Close()

	// 创建CSV writer
	csvWriter := csv.NewWriter(output)
	defer csvWriter.Flush()

	// 写入CSV头
	if err := csvWriter.Write([]string{"SPU", "Action", "ProductID"}); err != nil {
		return fmt.Errorf("写入CSV头失败: %v", err)
	}

	// 读取输入文件内容
	content, err := os.ReadFile("MK_AU.log")
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 按行处理
	lines := strings.Split(string(content), "\n")
	var records []Record

	for _, line := range lines {
		matches := logRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			records = append(records, Record{
				SPU:       matches[1],
				Action:    matches[2],
				ProductID: matches[3],
			})
		}
	}

	// 写入CSV数据
	for _, r := range records {
		row := []string{r.SPU, r.Action, r.ProductID}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("写入记录失败: %v", err)
		}
	}

	return nil
}

// 定义正则表达式匹配模式
var logRegex = regexp.MustCompile(`SPU:\s+(\S+).*?(保留|删除)商品:\s+(\S+)`)

// Record 表示解析后的记录
type Record struct {
	SPU       string
	Action    string
	ProductID string
}
