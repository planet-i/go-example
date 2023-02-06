package main

import (
	"fmt"
	"net/url"
)

func main() {
	uStr := "http://192.168.51.203:30648/v1/orgs/init-progress/1"
	u, _ := url.Parse(uStr)
	fmt.Printf("%v", u)
}
