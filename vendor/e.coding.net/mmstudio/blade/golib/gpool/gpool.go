package gpool

type Worker struct {
	JobChan chan Job

	// closeChan can be closed in order to cleanly shutdown this worker.
	closeChan chan struct{}
	// closedChan is closed by the run() goroutine when it exits.
	closedChan chan struct{}
}

func (w *Worker) Start(jobQueue chan Job) {
	defer func() {
		close(w.closedChan)
	}()

	go func() {
		var job Job
		for {
			select {
			case job = <-jobQueue:
				if job == nil {
					return
				}
				job()
			case <-w.closeChan:
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	close(w.closeChan)
}

func (w *Worker) join() {
	<-w.closedChan
}

func newWorker() *Worker {
	return &Worker{
		JobChan:    make(chan Job),
		closeChan:  make(chan struct{}),
		closedChan: make(chan struct{}),
	}
}

type Job func()

type Pool struct {
	JobQueue chan Job
	workers  []*Worker
}

func NewPool(numWorkers int, jobQueueLen int) *Pool {
	pool := &Pool{JobQueue: make(chan Job, jobQueueLen)}
	pool.SetSize(numWorkers)
	return pool
}

func (p *Pool) SetSize(n int) {
	lWorkers := len(p.workers)
	if lWorkers == n {
		return
	}

	// Add extra workers if N > len(workers)
	for i := lWorkers; i < n; i++ {
		w := newWorker()
		w.Start(p.JobQueue)
		p.workers = append(p.workers, w)
	}

	// Asynchronously stop all workers > N
	for i := n; i < lWorkers; i++ {
		p.workers[i].stop()
	}

	// Synchronously wait for all workers > N to stop
	for i := n; i < lWorkers; i++ {
		p.workers[i].join()
	}

	// Remove stopped workers from slice
	p.workers = p.workers[:n]

}

func (p *Pool) Release() {
	p.SetSize(0)
	close(p.JobQueue)
}
