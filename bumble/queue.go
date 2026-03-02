package bumble

import (
	"fmt"
	"sync"
)

type Job struct {
	id string
	fn func()
}

type Worker struct {
	id         int
	jobChannel chan Job
	wg         *sync.WaitGroup
}

func NewWorker(id int, jobChannel chan Job, wg *sync.WaitGroup) Worker {
	return Worker{
		id:         id,
		jobChannel: jobChannel,
		wg:         wg,
	}
}

func (w Worker) Start() {
	go func() {
		for job := range w.jobChannel {
			fmt.Printf("Worker %d started job %v\n", w.id, job.id)
			job.fn()
			fmt.Printf("Worker %d finished job %v\n", w.id, job.id)
			w.wg.Done()
		}
	}()
}

type Queue struct {
	jobChannel chan Job
	wg         *sync.WaitGroup
}

func NewQueue(numWorkers int) *Queue {
	jobChannel := make(chan Job)
	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, jobChannel, wg)
		worker.Start()
	}

	return &Queue{
		jobChannel: jobChannel,
		wg:         wg,
	}
}

func (q *Queue) Enqueue(label string, fn func()) {
	q.wg.Add(1)
	job := Job{
		id: label,
		fn: fn,
	}
	q.jobChannel <- job
}

func (q *Queue) Wait() {
	q.wg.Wait()
}

//func main() {
//	queue := NewQueue(3)
//
//	for i := 0; i < 5; i++ {
//		jobID := i + 1
//		job := Job{
//			id: jobID,
//			fn: func() {
//				time.Sleep(time.Second)
//				fmt.Printf("Job %d executed\n", jobID)
//			},
//		}
//		queue.Enqueue(job)
//	}
//
//	queue.Wait()
//}
