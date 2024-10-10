package main

import (
	"bytes"
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	excelizev2 "github.com/xuri/excelize/v2"
)

func main() {
	testImage()
	testExport()
}

func testExport() {
	f := excelize.NewFile()             //新建excel文件
	f.SetCellValue("Sheet1", "B2", 100) //填值
	f.SetCellValue("Sheet1", "A1", 90)
	f.SetCellValue("Sheet1", "c2", 60)
	f.SetCellValue("Sheet1", "c3", 150)

	//----增加一个工作表，添加图表
	index := f.NewSheet("sheet2")                                                                                              //添加新工作表
	categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"} //存行和列的说明
	values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}                  //存值
	for k, v := range categories {
		f.SetCellValue("sheet2", k, v)
	}
	for k, v := range values {
		f.SetCellValue("sheet2", k, v)
	}
	if err := f.AddChart("sheet2", "E1", `{"type":"col3DClustered","series":[{"name":"sheet2!$A$2","categories":"sheet2!$B$1:$D$1","values":"sheet2!$B$2:$D$2"},{"name":"sheet2!$A$3","categories":"sheet2!$B$1:$D$1","values":"sheet2!$B$3:$D$3"},{"name":"sheet2!$A$4","categories":"sheet2!$B$1:$D$1","values":"sheet2!$B$4:$D$4"}],"title":{"name":"Fruit 3D Clustered Column Chart"}}`); err != nil {
		fmt.Println(err)
		return
	}
	f.SetActiveSheet(index) //存储添加的工作表

	//----添加一个工作表，添加图片
	index2 := f.NewSheet("sheet3")
	if err := f.AddPicture("Sheet3", "A1", "image.png", ""); err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(index2)

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err) //存储建立的excel文件
	}

	f1, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	cell := f1.GetCellValue("sheet1", "B2")
	fmt.Println(cell)

	rows := f1.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

/* AddChart("当前工作表名", "图表左上角在表中的位置", `{"type":"图标类型",
"series":[{"name":"某类","categories":"x方向上的值的范围","values":"某类在x上对应的y值"},{"name":"Sheet1!$A$3","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$3:$D$3"},{"name":"Sheet1!$A$4","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$4:$D$4"}],
"title":{"name":"图表名称"}}`) */

func testImage() {
	// 创建一个新的Excel文件，并指定活跃的工作表
	activeSheet := "Sheet1"
	f := excelizev2.NewFile()
	index, _ := f.NewSheet(activeSheet)
	// 写表头
	colNos := make([]string, 0)
	colNos = append(colNos, []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"}...)
	headers := []string{"SPU", "图片"}
	for k, v := range headers {
		col := fmt.Sprintf("%v1", colNos[k])
		f.SetCellValue(activeSheet, col, v)
	}
	// 数据源
	res := [][]string{
		{"data1", "https://xiaowandou-ecom.s3.ap-southeast-1.amazonaws.com/upload/8b21d4e83432fe5d96638945baed3ee7.jpg"},
		{"data2", "https://cdn.unitemagic.com/upload/337c13d813fe46361f9b73f3558678a8.jpeg"},
		{"data3", "https://imgcdn.umcasual.com/um_ecom/66805433/upload/800d8a8ffc66878297772f389c733bf7_e772d785-0021-4521-bd55-1eae7c815adb.png"},
		{"data4", "https://cdn.unitemagic.com/creative/880567754520275/1705562178.jpeg"},
	}
	imageCol := GetImageCol(headers)
	// 插入数据和图片
	for i, row := range res {
		f.SetRowHeight(activeSheet, i+2, 60)
		for j, data := range row {
			col := fmt.Sprintf("%s%d", colNos[j], i+2) // 简化列号生成逻辑
			fmt.Println(col, data)
			if colNos[j] == imageCol && data != "" {
				locked := true
				format := &excelizev2.GraphicOptions{
					LockAspectRatio: true,
					Locked:          &locked,
					Positioning:     "oneCell",
					AutoFit:         true,
				}

				// 读取远程图片
				image, err := ReadRemoteFile(data)
				if err != nil {
					log.Printf("ReadRemoteFile err: %s", err)
				}
				ext := GetImageExt(data)
				fmt.Println(ext)
				if err := f.AddPictureFromBytes(activeSheet, col, "image", ext, image, format); err != nil {
					log.Printf("f.AddPictureFromBytes err: %s", err.Error())
				}
				f.SetCellHyperLink(activeSheet, col, data, "External")
				continue
			}
			f.SetCellValue(activeSheet, col, data)
		}
	}
	f.SetActiveSheet(index)
	// 保存文件
	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatalf("Save file failed, %s", err)
	}
	fmt.Println("Excel file created successfully.")
}

// 读取远程图片
func ReadRemoteFile(url string) ([]byte, error) {
	result := make([]byte, 0)
	if url == "" {
		return result, errors.New("URL不能为空")
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	buffer := bytes.NewBuffer(make([]byte, 0, 2048))
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GetImageCol(headers []string) string {
	index := -1
	for i, header := range headers {
		if header == "图片" {
			index = i
		}
	}
	if index > -1 {
		col, _ := excelizev2.ColumnNumberToName(index + 1)
		return col
	}
	return ""
}

func GetImageExt(path string) string {
	// 查找最后一个 "." 的位置
	lastDotIndex := strings.LastIndex(path, ".")

	// 如果没有找到 "." 或者 "." 在字符串开头，则返回空字符串
	if lastDotIndex == -1 || lastDotIndex == 0 {
		return ".jpg"
	}

	// 提取 "." 之后的部分
	ext := path[lastDotIndex:]
	return strings.ToLower(ext)
}
