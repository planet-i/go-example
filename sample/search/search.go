package search

import (
	"log"
	"sync"
)

//匹配器映射
var matchers = make(map[string]Matcher)

func Run(searchTerm string) {
	feeds, err := RetrieveFeeds() //获取需要搜索的数据源列表
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan *Result)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(feeds))

	for _, feed := range feeds {
		matcher, exists := matchers[feed.Type] //获取匹配器
		if !exists {
			matcher = matchers["default"]
		}
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			//匹配器，数据源，搜索项、匹配结果输出懂啊result通道
			waitGroup.Done()
		}(matcher, feed)
	}
	go func() {
		waitGroup.Wait()
		close(results) //关闭通道
	}()
	Display(results)
}
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
