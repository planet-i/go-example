package main

import (
	"fmt"
	"time"

	"math/rand"
)

type prize struct {
	Id         int64   // 奖品ID
	ActiveId   string  // 活动ID
	Class      int64   // 奖品类别
	Name       string  // 奖品名称
	TotalNum   int64   // 奖品总数量
	CurrentNum int64   // 奖品当前数量
	Ratio      float64 // 奖品中奖概率
	Created    int64   // 创建时间
	Updated    int64   // 更新时间
}

func main() {
	// 示例奖品列表
	prizes := []*prize{
		{Id: 1, ActiveId: "6cfc2a51", Class: 1, Name: "一等奖", Ratio: 0.06, TotalNum: 20, CurrentNum: 20},
		{Id: 2, ActiveId: "6cfc2a51", Class: 1, Name: "二等奖", Ratio: 0.16, TotalNum: 50, CurrentNum: 50},
		{Id: 3, ActiveId: "6cfc2a51", Class: 1, Name: "三等奖", Ratio: 0.31, TotalNum: 100, CurrentNum: 100},
		{Id: 4, ActiveId: "6cfc2a51", Class: 2, Name: "四等奖", Ratio: 0.47, TotalNum: 150, CurrentNum: 150},
	}

	// 进行100次抽奖
	var one, two, three, four int
	for i := 0; i < 10000; i++ {
		p := drawPrize(prizes)
		if p == nil {
			fmt.Printf("第 %d 次抽奖: 没有中奖\n", i+1)
			continue
		}
		fmt.Printf("第 %d 次抽奖: 恭喜您获得：%s\n", i+1, p.Name)
		switch p.Name {
		case "一等奖":
			one += 1
		case "二等奖":
			two += 1
		case "三等奖":
			three += 1
		case "四等奖":
			four += 1
			code := GenerateDiscountCode(12)
			fmt.Printf("优惠码:%s\n", code)
		}
	}
	fmt.Printf("一等奖发出个数:%d\n", one)
	fmt.Printf("二等奖发出个数:%d\n", two)
	fmt.Printf("三等奖发出个数:%d\n", three)
	fmt.Printf("四等奖发出个数:%d\n", four)
}

// GetRandomInt 返回一个在 [1,m] 区间的随机整数
func GetRandomInt(m int) int {
	if m <= 1 {
		return m
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(m) + 1
}

// 抽奖函数。规则：按照设置的奖品概率抽奖。抽到四等奖不更改库存。如果抽到无库存的奖品，视为中四等奖
func drawPrize(prizes []*prize) *prize {
	num := len(prizes)
	if num == 0 {
		return nil
	}
	random := GetRandomInt(100)
	ratio := 0
	for _, p := range prizes {
		ratio += int(p.Ratio * 100)
		fmt.Println("随机数比较", random, ratio)
		if random <= ratio {
			if p.CurrentNum > 0 {
				if p.Class == 1 {
					p.CurrentNum--
				}
				return p
			}
			fmt.Println("商品库存不够")
			return nil
		}
	}
	return nil
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateDiscountCode 生成一个指定长度的随机折扣码
func GenerateDiscountCode(length int) string {
	if length <= 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
