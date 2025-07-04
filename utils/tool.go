package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mozillazg/go-pinyin"
)

// GetFilenameFromUrl 函数从给定的 URL 中提取文件名和不带扩展名的文件名
// 返回值为文件名和不带扩展名的文件名
func GetFilenameFromUrl(url string) (fullName, name string) {
	arr := strings.Split(url, "/")
	if len(arr) >= 1 {
		fullName = arr[len(arr)-1]
		name = strings.Split(fullName, ".")[0]
	}
	return
}

// 获取字符串的第一个字
func firstLetter(s string) string {
	_, size := utf8.DecodeRuneInString(s)
	return s[:size]
}

func firstLetter1(s string) string {
	for _, l := range s {
		return string(l)
	}
	return ""
}

func FirstLetterOfPinYin(r string) string {
	var a = pinyin.NewArgs()
	result := pinyin.Pinyin(r, a)
	fmt.Println(result)
	return string(result[0][0][0])
}

// 获取字符串首字母
func FirstLetterOfString(s string) string {
	l := []rune(s)[:1][0]
	if unicode.Is(unicode.Han, l) {
		fmt.Println("是一个汉字", string(l))
		var a = pinyin.NewArgs()
		result := pinyin.Pinyin(string(l), a)
		return string(result[0][0][0])
	} else if unicode.IsNumber(l) {
		fmt.Println("是一个数字", string(l))
		return "~"
	} else if unicode.IsLetter(l) {
		fmt.Println("是一个字母", string(l))
		a := string(unicode.ToUpper(l))
		b := string(unicode.ToLower(l))
		fmt.Println(string(a), string(b))
		return a
	}
	return ""
}

func GetExt() {
	var a = "gggg.pdf.txt"
	fmt.Println(strings.TrimPrefix(a, "."))
}

type prize struct {
	Id           int64
	ActiveId     string
	ActiveSiteId []int64
	Level        int64   // 奖品等级
	Name         string  // 奖品名称
	TotalNum     int64   // 奖品总数量
	Ratio        float64 // 中奖概率
	CurrentNum   int64   // 当前的数量
	Created      int64
	Updated      int64
}

func drawPrize(prizes []prize) (prizell prize, ok bool) {
	num := len(prizes)
	if num == 0 {
		return prize{}, false
	}
	random := GetRandomInt(100)
	// 如果中奖的级别是1等奖且奖品的当前数量>0,就中1等奖
	// 如果中奖的级别是2等奖且奖品的当前数量>0,就中2等奖
	// 如果中奖的级别是3等奖且奖品的当前数量>0,就中3等奖
	// 如果中奖的级别是4等奖且奖品的当前数量>0,就中4等奖
	// 否则返回没中奖
	fmt.Println("随机数", random)
	ratio := 0
	for _, p := range prizes {
		ratio += int(p.Ratio * 100)
		fmt.Println("随机数比较", random, ratio)
		if random <= ratio && p.CurrentNum > 0 {
			fmt.Println("中奖啦", random, ratio, p.Level, p.CurrentNum)
			return p, true
		}
	}
	return prize{}, false
}

// GetRandomInt 返回一个在 [1,m] 区间的随机整数
func GetRandomInt(m int) int {
	if m <= 1 {
		return m
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(m) + 1
}

func testPrize() {
	// 示例奖品列表
	prizes := []prize{
		{Id: 1, Level: 1, Name: "一等奖", Ratio: 0.1, CurrentNum: 5},
		{Id: 2, Level: 2, Name: "二等奖", Ratio: 0.2, CurrentNum: 10},
		{Id: 3, Level: 3, Name: "三等奖", Ratio: 0.3, CurrentNum: 15},
		{Id: 4, Level: 4, Name: "四等奖", Ratio: 0.4, CurrentNum: 20},
	}

	// 进行一次抽奖
	p, ok := drawPrize(prizes)
	if ok {
		fmt.Printf("恭喜您获得：%s\n", p.Name)
	} else {
		fmt.Println("很遗憾，您没有中奖")
	}
}
