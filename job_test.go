package wpool

import (
	"log"
	"testing"
)

func TestJob(t *testing.T) {
	log.Print("--------------------- Testing job")
	testJob := NewJob("test-job", func(job JobInterface) (err error) {
		log.Print(job.Name() + " in progress")
		return
	})
	log.Print("run " + testJob.Name())
	err := testJob.Run()
	log.Printf("%s error: %v", testJob.Name(), testJob.Error())
	log.Printf("%s duration: %s", testJob.Name(), testJob.Duration())
	if err != nil {
		t.Error(err)
	}
}
