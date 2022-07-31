package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

//管理资源的池
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer //被管理的资源需要实现io.Closer接口
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("Pool has been closed.")

//规定池的大小、分配资源的函数
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value too small")
	}
	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}
func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		r.Close() //想放回去的时候，池已经关闭了，就销毁这个资源
		return
	}
	select {
	case p.resources <- r: //把资源r放回资源池
		log.Println("Release:", "In Queue")
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock() //44 45  59  60  是对同一个互斥量进行加锁和解锁，阻止了两个方法在不同的goroutine里同时运行

	if p.closed {
		return
	}
	p.closed = true //将池关闭
	close(p.resources)
	for r := range p.resources {
		r.Close()
	} //
}
