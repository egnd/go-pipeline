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

type JobFunction func(job JobInterface) error

type Job struct {
	JobInterface
	name     string
	err      error
	duration time.Duration
	callback JobFunction
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
	this.err = this.logic()
	this.duration = time.Since(start)
	return this.Error()
}

func (this *Job) logic() error {
	return this.callback(this)
}

func NewJob(name string, jobLogic JobFunction) JobInterface {
	return &Job{name: name, callback: jobLogic}
}
