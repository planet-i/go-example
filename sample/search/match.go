package search

import (
	"fmt"
	"log"
)

type Result struct {
	Field   string
	Content string
}
type Matcher interface {
	Search(feed *Feed, seachTerm string) ([]*Result, error)
}

func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	//对特定的匹配器执行搜索
	searchResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}
	for _, result := range searchResults {
		results <- result //写入结果
	} //在search.Run中close通道
}

func Display(results chan *Result) {
	for result := range results {
		fmt.Printf("Display %s:\n%s\n\n", result.Field, result.Content)
	}
} //通道会一直被阻塞，直到有结果写入
