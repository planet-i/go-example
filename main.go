package main

import (
	"os"
	"runtime/trace"
)

func main() {
	testAllocate()
	testTraceTool()
}

func allocate() {
	_ = make([]byte, 1<<20)
}

func testAllocate() {
	for i := 1; i < 10000; i++ {
		allocate()
	}
} // 1. go build -o main   2. GODEBUG=gctrace=1 ./main

var cache = map[interface{}]interface{}{}

func keepalloc() {
	for i := 0; i < 10000; i++ {
		m := make([]byte, 1<<10)
		cache[i] = m
	}
}

func keepalloc2() {
	for i := 0; i < 10000; i++ {
		go func() {
			select {}
		}()
	}
}

var ch = make(chan struct{})

func keepalloc3() {
	for i := 0; i < 10000; i++ {
		go func() {
			ch <- struct{}{}
		}()
	}
}
func testTraceTool() {
	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	keepalloc()
	keepalloc2()
	keepalloc3()
} // 1. go run main.go    2. go tool trace trace.out
