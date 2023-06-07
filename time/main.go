// package main

// import (
// 	"fmt"
// 	"strconv"
// 	"time"
// )

// func main() {
// 	//Trans()
// 	//aa()
// 	bb()
// 	// const TimeLayout = "2006-01-02 15:04:05"
// 	// fmt.Println(time.Now())
// 	// fmt.Println(time.Now().Unix()) //将t表示为Unix时间，即从时间点January 1, 1970 UTC到时间点t所经过的时间（单位秒）。
// 	// fmt.Println(time.Now().Format(TimeLayout))
// 	// timeStr := "2017-11-13 11:14:21"

// 	// loc, _ := time.LoadLocation("Local")
// 	// tm, error := time.ParseInLocation(TimeLayout, strings.Trim(timeStr, "\n"), loc)
// 	// if error != nil {
// 	// 	log.Fatal(error)
// 	// }
// 	// fmt.Println(tm.Unix())

// 	// const layout = "Jan 2, 2006 at 3:04pm (MST)"
// 	// t := time.Date(2009, time.November, 10, 15, 0, 0, 0, time.Local)
// 	// fmt.Println(t.Format(layout))
// 	// fmt.Println(t.UTC().Format(layout))
// }

// func Trans() {
// 	str := "2022-08-31T09:36:14.934041+08:00"
// 	t1, err := time.Parse("2006-01-02 15:04:05", str)
// 	fmt.Println(t1, err)
// }
// func aa() {
// 	fmt.Println(time.Now().Unix())
// 	b := strconv.FormatInt(time.Now().Unix(), 10)
// 	fmt.Println(b)
// }

// func bb() {
// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
// 	fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"))
// }

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type MyStruct struct {
	Name  string
	Time  time.Time
	DTime *time.Time
}

func main() {
	jsonString := `{"Name": "John", "Time": "2023-02-13T08:54:19.190735Z","DTime":null}`
	var myStruct MyStruct
	err := json.Unmarshal([]byte(jsonString), &myStruct)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(myStruct.Name, myStruct.Time, myStruct.DTime)
}
