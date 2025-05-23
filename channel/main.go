package main

import (
	"errors"
	"fmt"
)

func main() {
	ch := make(chan BatchData, 1)
	go info(ch)
	consumer(ch)
}

type BatchData struct {
	Data        []int // 当前批次数据
	Total       int   // 总数据量
	BatchSeq    int   // 批次序号
	IsLastBatch bool  // 是否是最后一批
}

func info(ch chan BatchData) (total int64, err error) {
	page := 1
	pageSize := 10
	for {
		rows, err := GetData(page, pageSize)
		if err != nil {
			return total, err
		}
		batch := rows.Data
		data := BatchData{
			Data:        batch,
			Total:       rows.TotalCount,
			BatchSeq:    page,
			IsLastBatch: !rows.HasNext,
		}
		ch <- data // 发送批次数据到Channel
		if !rows.HasNext {
			break
		}
		page++
	}
	close(ch) // 数据获取完成后关闭Channel
	return
}

func consumer(ch chan BatchData) {
	var i = 1
	for batch := range ch {
		fmt.Printf("获取到第 %d 批:最后一批？%v\n", batch.BatchSeq, batch.IsLastBatch)
		for _, v := range batch.Data {
			println(v)
		}
		i++
	}
}

type DataInfo struct {
	TotalCount int
	Data       []int
	HasNext    bool
}

func GetData(page, pageSize int) (*DataInfo, error) {
	total := 100
	// 参数校验
	if page < 1 {
		return nil, errors.New("页码不能小于1")
	}
	if pageSize < 1 {
		return nil, errors.New("每页数量不能小于1")
	}

	// 计算起始位置
	start := (page - 1) * pageSize
	if start >= total {
		return &DataInfo{
			TotalCount: total,
			Data:       []int{},
			HasNext:    false,
		}, nil
	}

	// 计算结束位置
	end := start + pageSize
	if end > total {
		end = total
	}

	// 生成数据
	data := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		data = append(data, i)
	}

	// 判断是否有下一页
	hasNext := end < total

	return &DataInfo{
		TotalCount: total,
		Data:       data,
		HasNext:    hasNext,
	}, nil
}
