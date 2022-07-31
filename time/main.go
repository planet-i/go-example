package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	const TimeLayout = "2006-01-02 15:04:05"
	fmt.Println(time.Now())
	fmt.Println(time.Now().Unix()) //将t表示为Unix时间，即从时间点January 1, 1970 UTC到时间点t所经过的时间（单位秒）。
	fmt.Println(time.Now().Format(TimeLayout))
	timeStr := "2017-11-13 11:14:21"

	loc, _ := time.LoadLocation("Local")
	tm, error := time.ParseInLocation(TimeLayout, strings.Trim(timeStr, "\n"), loc)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(tm.Unix())

	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	t := time.Date(2009, time.November, 10, 15, 0, 0, 0, time.Local)
	fmt.Println(t.Format(layout))
	fmt.Println(t.UTC().Format(layout))

}
