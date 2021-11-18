package pool

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	RUNNING = iota
	STOPPING
	STOPED
)

var (
	ErrPoolNotRunning  = errors.New("the pool is not running")
	ErrInvalidPoolSize = errors.New("invalid pool size")
)

type Worker struct {
	Process func(param ...interface{}) interface{}
	Result  interface{}
	Param   []interface{}
}

type Pool struct {
	size           int
	status         int
	runningWorkers uint64
	resultCh       chan *Worker
	workerCh       chan *Worker
	workerCache    []*Worker
	sync.Mutex
}

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrInvalidPoolSize
	}
	p := &Pool{
		size:        size,
		status:      RUNNING,
		resultCh:    make(chan *Worker),
		workerCh:    make(chan *Worker, size),
		workerCache: make([]*Worker, size),
	}
	return p, nil
}

func (p *Pool) Put(worker *Worker) error {
	if p.getStatus() != RUNNING {
		return ErrPoolNotRunning
	}
	p.workerCache = append(p.workerCache, worker)
	return nil
}

func (p *Pool) Run() <-chan *Worker {
	if p.getStatus() != RUNNING {
		return nil
	}
	p.run()
	return p.resultCh
}

func (p *Pool) run() {
	go func() {
		for {
			w := p.getWorker()
			if w != nil {
				p.workerCh <- w
			}
			if len(p.workerCache) == 0 {
				close(p.workerCh)
				p.setStatus(STOPPING)
				return
			}
		}
	}()

	go func() {
		for {
			if p.running() == 0 {
				close(p.resultCh)
				p.setStatus(STOPED)
				return
			}
		}
	}()

	for i := p.running(); i < p.size; i++ {
		p.incRunning()
		go func() {
			defer p.decRunning()
			for worker := range p.workerCh {
				worker.Result = worker.Process(worker.Param...)
				p.resultCh <- worker
			}
		}()
	}
}

func (p *Pool) getWorker() *Worker {
	p.Lock()
	defer p.Unlock()
	if len(p.workerCache) > 0 {
		w := p.workerCache[0]
		p.workerCache = p.workerCache[1:]
		return w
	}
	return nil
}

func (p *Pool) getStatus() int {
	p.Lock()
	defer p.Unlock()
	return p.status
}

func (p *Pool) setStatus(status int) {
	p.Lock()
	defer p.Unlock()
	p.status = status
}

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

func (p *Pool) running() int {
	return int(atomic.LoadUint64(&p.runningWorkers))
}
