package work

import "sync"

type Worker interface {
	Task()
} //类比io.Closer角色

//工作池
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

//工作池只会包含 maxGoroutines个执行任务的goroutine
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work { //阻塞等待 ，直到Run通过p.work  传值过来
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return &p
}

//Run提交工作到池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

//等待所有goroutine停止工作
func (p *Pool) Shutdown() {
	close(p.work) //关闭通道，导致所有池里的goroutine停止工作
	p.wg.Wait()
}
