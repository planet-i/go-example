package main

import (
	"log"
	"os"

	_ "github.com/planet-i/go-example/sample/matchers" //注册RSS匹配器
	"github.com/planet-i/go-example/sample/search"     //init函数注册默认匹配器
)

func init() {
	//日志输出：标准错误stderr  -->  标准输出stdout
	log.SetOutput(os.Stdout)
}
func main() {
	search.Run("president")
	//time.Sleep(1 * time.Second)
}
