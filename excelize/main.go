package main

import (
	"fmt"
	_ "image/png"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
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
