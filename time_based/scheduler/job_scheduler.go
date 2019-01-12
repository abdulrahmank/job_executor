package scheduler

import (
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"log"
)

type TimeBasedScheduler struct {
	Executor exector.Executor
}

func (t *TimeBasedScheduler) Schedule(time, filename string) {
	_, e := t.Executor.ExecuteJob(filename)
	if e != nil {
		log.Fatalf("Unable to execute %s due to %v\n", filename, e)
	}
}
