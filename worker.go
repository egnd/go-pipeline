package wpool

import "log"

type WorkerInterface interface {
	Listen(pipeline chan chan JobInterface)
	Stop()
}

type Worker struct {
	WorkerInterface
	pipe chan JobInterface
	stop chan bool
}

func (this *Worker) Listen(pipeline chan chan JobInterface) {
	this.stop = make(chan bool)
	this.pipe = make(chan JobInterface)
	go func() {
		log.Println("worker start")
		for {
			select {
			case <-this.stop:
				log.Println("worker stop")
				return
			case job := <-this.pipe:
				log.Println("worker execute job " + job.Name())
				this.runJob(job)
			default:
				log.Println("worker waiting")
				pipeline <- this.pipe
			}
		}
	}()
}

func (this *Worker) Stop() {
	this.stop <- true
	close(this.pipe)
	close(this.stop)
}

func (this *Worker) runJob(job JobInterface) {
	job.Run()
}
