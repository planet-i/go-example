package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func main() {
	splitPDF()  // 分割PDF文件
	mergeFile() // 合并PDF文件
	getImage()  // 把PDF转成图片
}

// mergeFile 合并多个PDF文件
func mergeFile() {
	// 合并文件
	inDir := "."
	outDir := "."
	inFile1 := filepath.Join(inDir, "1.pdf")
	inFile2 := filepath.Join(inDir, "2.pdf")
	inFile3 := filepath.Join(inDir, "3.pdf")
	outFile := filepath.Join(outDir, "merge.pdf")
	mergePDFs([]string{inFile1, inFile2, inFile3}, outFile)
}

func mergePDFs(inputFiles []string, outputPath string) error {
	if err := api.MergeCreateFile(inputFiles, outputPath, false, nil); err != nil {
		return fmt.Errorf("PDF合并失败: %v", err)
	}
	return nil
}

// getImage 将PDF转换为图片
func getImage() { //
	inputPDF := "1.pdf"
	outputImage := "output.jpg"

	// 调用Ghostscript转图片（需安装）
	cmd := exec.Command("gs",
		"-sDEVICE=pngalpha",
		"-r300",                // 分辨率300 DPI
		"-dNOPAUSE", "-dBATCH", // 无交互模式
		"-sOutputFile="+outputImage,
		inputPDF,
	)
	if err := cmd.Run(); err != nil {
		fmt.Printf("图片转换失败: %v\n", err)
		return
	}
}

// splitPDF 拆分PDF文件
func splitPDF() {
	inputFile := "merge.pdf" // 输入PDF文件路径（路径处理要注意）
	outputDir := "."         // 输出目录

	// 1. 按单页拆分（每页生成单独文件）
	if err := splitBySinglePage(inputFile, outputDir); err != nil {
		fmt.Printf("拆分失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PDF已按单页拆分完成")

	// 2. 按指定页面范围拆分
	pageRange := "3-5" // 提取第3-5页
	if err := extractPageRange(inputFile, outputDir+"/range", pageRange); err != nil {
		fmt.Printf("范围提取失败: %v\n", err)
	} else {
		fmt.Println("指定页面范围已提取")
	}

	// 3. 按固定页数拆分（每2页一个文件）
	if err := splitByPageCount(inputFile, outputDir+"/groups", 2); err != nil {
		fmt.Printf("分组拆分失败: %v\n", err)
	} else {
		fmt.Println("按页数分组拆分完成")
	}
}

// 按单页拆分PDF
func splitBySinglePage(inputFile, outputDir string) error {
	return api.SplitFile(inputFile, outputDir, 1, nil) // 参数1表示每页拆分
}

// 提取指定页面范围
func extractPageRange(inputFile, outputDir, pageRange string) error {
	return api.ExtractPagesFile(inputFile, outputDir, []string{pageRange}, nil)
}

// 按固定页数分组拆分
func splitByPageCount(inputFile, outputDir string, pageCount int) error {
	return api.SplitFile(inputFile, outputDir, pageCount, nil)
}
