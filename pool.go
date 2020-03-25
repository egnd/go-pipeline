package wpool

import "log"

type PoolInterface interface {
	AddWorker(worker WorkerInterface) PoolInterface
	AddJob(job JobInterface) PoolInterface
	Stop()
}

type Pool struct {
	PoolInterface
	workers  []WorkerInterface
	pipeline chan chan JobInterface
	jobs     chan JobInterface
	stop     chan bool
}

func (this *Pool) AddWorker(worker WorkerInterface) PoolInterface {
	log.Println("pool add worker")
	worker.Listen(this.pipeline)
	this.workers = append(this.workers, worker)
	return this
}

func (this *Pool) AddJob(job JobInterface) PoolInterface {
	log.Println("pool add job " + job.Name())
	this.jobs <- job
	return this
}

func (this *Pool) Stop() {
	for _, worker := range this.workers {
		worker.Stop()
	}
	this.stop <- true
	close(this.stop)
	close(this.pipeline)
	close(this.jobs)
}

func NewPool() PoolInterface {
	workersPipe := make(chan chan JobInterface)
	jobsPipe := make(chan JobInterface)
	stop := make(chan bool)
	go func() {
		log.Println("pool start")
		var job JobInterface
		jobs := []JobInterface{}
		for {
			select {
			case <-stop:
				log.Println("pool stop")
				return
			case job = <-jobsPipe:
				log.Println("pool recieve job")
				jobs = append(jobs, job)
			case workerPipe := <-workersPipe:
				log.Println("pool recieve waiting worker")
				if len(jobs) > 0 {
					job, jobs = jobs[0], jobs[1:]
					log.Println("pool send to worker job " + job.Name())
					workerPipe <- job
				}
			default:
			}
		}
	}()
	return &Pool{pipeline: workersPipe, jobs: jobsPipe, stop: stop}
}
