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
