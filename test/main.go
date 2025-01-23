package main

import (
	"fmt"
	"io"
)

// Product 表示商品的价格信息
type Product struct {
	ProductId int64
	Price     float64
}

func calculatePrices(stream io.Reader) (minPrice, maxPrice, price, compareAtPrice float64, count int, err error) {
	// 初始化变量
	minPrice = 0
	maxPrice = 0
	price = 0
	compareAtPrice = 0
	count = 0

	// 假设 stream 是一个可以接收数据的流
	for {
		// 假设通过 stream 读取数据（这里只是一个模拟，实际上应该是从流中解析数据）
		// 这里我们需要模拟 stream.Recv() 的行为，通常是处理接收到的 Product 数据
		var pv Product      // 假设流返回的是 Product 对象
		err = mockRecv(&pv) // 这里 mockRecv 会模拟接收一个 Product，实际应是流的接收方法
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("stream.Recv(): %s", err)
		}

		// 计算最小价格
		if pv.Price < minPrice || minPrice == 0 {
			minPrice = pv.Price
		}

		// 计算最大价格
		if pv.Price > maxPrice {
			maxPrice = pv.Price
		}

		// 记录第一个商品价格
		if price == 0 {
			price = pv.Price
		}

		// 记录比较价格
		if compareAtPrice == 0 {
			compareAtPrice = pv.Price
		}

		count++
	}

	return minPrice, maxPrice, price, compareAtPrice, count, nil
}

// 模拟从流中接收数据（这里用来测试，实际情况会从 stream 中获取数据）
func mockRecv(pv *Product) error {
	// 这里模拟返回数据流，你可以在这里修改测试数据
	if *pv == (Product{}) {
		*pv = Product{ProductId: 1, Price: 10} // 第一个商品价格为 10
	} else if pv.ProductId == 1 {
		*pv = Product{ProductId: 2, Price: 20} // 第二个商品价格为 20
	} else if pv.ProductId == 2 {
		*pv = Product{ProductId: 3, Price: 5} // 第三个商品价格为 5
	} else {
		return io.EOF // 假设数据流结束
	}
	return nil
}

func main() {
	// 模拟数据流的处理
	minPrice, maxPrice, price, compareAtPrice, count, err := calculatePrices(nil)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印计算结果
	fmt.Printf("Total Products: %d\n", count)
	fmt.Printf("Min Price: %.2f\n", minPrice)
	fmt.Printf("Max Price: %.2f\n", maxPrice)
	fmt.Printf("Price (First Product): %.2f\n", price)
	fmt.Printf("Compare At Price (First Product): %.2f\n", compareAtPrice)
}
