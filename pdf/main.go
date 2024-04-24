package main

import (
	"fmt"

	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// 创建新的 PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// 添加需要合并的 PDF 文件
	for i := 1; i < 43; i++ {
		filepath := "C:\\Users\\EDY\\Downloads\\111" + "1 (" + strconv.Itoa(i) + ").pdf"
		pdf.AddPage()
		pdf.SetSourceFile(file)
		tplIdx := pdf.ImportPage(1)
		pdf.UseTemplate(tplIdx, 10, 10, 200, 0)
	}

	// 保存合并后的 PDF 文件
	err := pdf.OutputFileAndClose("merged_file.pdf")
	if err != nil {
		fmt.Println("Error saving PDF file:", err)
		return
	}

	fmt.Println("PDF files merged successfully.")
}
