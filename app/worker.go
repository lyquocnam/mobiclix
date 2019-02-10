package app

type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
}

func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
	}
}

func (w *Worker) Start() {

}

func (w *Worker) Stop() {

}
