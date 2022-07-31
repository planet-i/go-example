package main

import "github.com/jinzhu/gorm"

type Test struct {
	ID   int64  `gorm:"type:bigint(20);column:id;primary_key`
	Name string `gorm:"type:varchar(20);column:name`
	Age  int    `gorm:"type:int(11);column:age`
}

//每个字段后面的gorm是结构标记，可以用于声明对应数据库字段的属性。
func main() {
	db, connerr := gorm.Open("mysql", "root:1234567890@(127.0.0.1:3306)/dbf?charset=utf8&parseTime=True&loc=Local")
	if connerr != nil {
		panic("falied to connect database")
	}
	defer db.Close()
	db.SingularTable(true)

	test := &Test{
		ID:   3,
		Name: "Jkson",
		Age:  12,
	}
	db.Create(test)
	db.Model(&test).Update("name", "Linda")
	var testResult Test
	db.Where("name = ?", "Linda").First(&testResult)
	db.Delete(test)
}
