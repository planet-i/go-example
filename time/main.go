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
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type MyStruct struct {
	Name  string
	Time  time.Time
	DTime *time.Time
}

func main() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	nowTime := time.Now()
	locTime := nowTime.In(loc)
	fmt.Println(nowTime)
	fmt.Println(locTime)
	fmt.Println(nowTime.Format("2006-01-02 15:04:05"))
	fmt.Println(locTime.Format("2006-01-02 15:04:05"))
	fmt.Println(nowTime.Unix())
	fmt.Println(locTime.Unix())
	//d()
	//a()
	//b()
	//c()
}

func d() {
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, "2024-10-07 00:00:00")
	timestamp := t.Unix() // 获取 Unix 时间戳
	fmt.Println(timestamp)
}

func a() {
	begin := time.Now().AddDate(0, 0, -7).Format("20060102")
	beginInt, err := strconv.ParseInt(begin, 10, 64)
	if err != nil {
		logrus.Debugf(" strconv.ParseInt(begin, 10, 64): %s", err)
	}
	end := time.Now().Format("20060102")
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		logrus.Debugf(" strconv.ParseInt(end, 10, 64): %s", err)
	}
	fmt.Println(beginInt, endInt)
}

func b() {
	var enhour int64 = 1723098972
	m := time.Unix(enhour, 0).Format("2006-01-02 15:04:05")
	fmt.Println(m)
}

func c() {
	str := "2024082122"
	layout := "2006010215"
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println(t.Format("2006-01-02 15")) // 格式化为易读时间
}
