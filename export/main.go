package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	activeSheet := "Sheet1"
	xfile := excelize.NewFile()           // 使用新的库创建文件
	_, err := xfile.NewSheet(activeSheet) // 使用新的库创建工作表
	if err != nil {
		fmt.Println(err.Error())
	}
	s, err := excelize.ColumnNumberToName(13)
	fmt.Println(s, err)
}

// 	// 添加表头
// 	headers := []string{"订单", "订单ID", "下单时间", "产品规格", "原始产品SKU", "产品SKU", "尺寸", "数量", "定制方案", "改单时间", "状态", "原因", "图片"}
// 	for i, cellData := range headers {
// 		col, err := excelize.ColumnNumberToName(i + 1)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		cell := fmt.Sprintf("%v1", col)
// 		//cell := fmt.Sprintf("A%d", i+1) // 根据列索引计算单元格位置
// 		fmt.Println(cell, cellData)
// 		err = xfile.SetCellValue(activeSheet, cell, cellData) // 使用新的库设置单元格值
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 	}
// 	cellData := "https://img-va.myshopline.com/image/store/1715325058293/01.jpeg?w=794\\u0026h=851"
// 	file, _ := ReadRemoteFile(cellData)
// 	img, err := jpeg.Decode(bytes.NewReader(file))
// 	if err != nil {
// 		fmt.Printf("jpeg.Decode err: %s", err.Error())
// 	}
// 	encodebuffer := bytes.NewBuffer(nil)
// 	jpeg.Encode(encodebuffer, img, nil)
// 	locked := true
// 	format := &excelize.GraphicOptions{
// 		LockAspectRatio: true,
// 		Locked:          &locked,
// 		Positioning:     "oneCell",
// 		AutoFit:         true,
// 	}
// 	xfile.AddPictureFromBytes(activeSheet, "A2", "image", ".jpg", encodebuffer.Bytes(), format)
// 	xfile.SetCellHyperLink(activeSheet, "A2", cellData, "External")
// }

// func main() {
// 	f := excelize.NewFile()
// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	file, err := os.ReadFile("image.jpg")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	if err := f.AddPictureFromBytes("Sheet1", "A2", &excelize.Picture{
// 		Extension: ".jpg",
// 		File:      file,
// 		Format:    &excelize.GraphicOptions{AltText: "Excel Logo"},
// 	}); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	if err := f.SaveAs("Book1.xlsx"); err != nil {
// 		fmt.Println(err)
// 	}
// }

// func ReadRemoteFile(url string) ([]byte, error) {
// 	result := make([]byte, 0)
// 	if url == "" {
// 		return result, errors.New("URL不能为空")
// 	}

// 	var netClient = &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	resp, err := netClient.Get(url)
// 	if err != nil {
// 		return result, err
// 	}
// 	defer resp.Body.Close()

// 	result, _ = io.ReadAll(resp.Body)
// 	return result, nil
// }
