package app

type Dispatcher struct {
	maxWorkers int
	jobQueue   chan Job
	workerPool chan chan Job
}

func (d *Dispatcher) NewDispatcher(maxWorkers int, jobQueue chan Job) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		workerPool: make(chan chan Job),
		jobQueue:   jobQueue,
	}
}
