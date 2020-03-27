package wpool

import (
	"time"
)

type JobInterface interface {
	Name() string
	Error() error
	Duration() time.Duration
	Run() error
}

type JobCallback func(job JobInterface) error

type Job struct {
	JobInterface
	name     string
	err      error
	duration time.Duration
	callback JobCallback
}

func (this *Job) Name() string {
	return this.name
}

func (this *Job) Error() error {
	return this.err
}

func (this *Job) Duration() time.Duration {
	return this.duration
}

func (this *Job) Run() (err error) {
	start := time.Now()
	defer func() {
		this.duration = time.Since(start)
	}()
	this.err = this.callback(this)
	return this.Error()
}

func NewJob(name string, jobLogic JobCallback) JobInterface {
	return &Job{name: name, callback: jobLogic}
}
