package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numberGoroutines = 4
	taskLoad         = 10
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().Unix())
}
func main1() {
	tasks := make(chan string, taskLoad)
	wg.Add(numberGoroutines)

	for gr := 1; gr <= numberGoroutines; gr++ {
		go worker(tasks, gr)
	}
	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("Task: %d", post)
	}
	close(tasks)
	wg.Wait()
}
func worker(tasks chan string, worker int) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			fmt.Printf("Worker:%d:Shutting Down\n", worker)
			return
		}
		fmt.Printf("worker :%d :Started %s\n", worker, task)

		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Microsecond)
		fmt.Printf("worker : %d:Complated %s\n", worker, task)
	}
}
