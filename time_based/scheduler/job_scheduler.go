package scheduler

import (
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"log"
	"time"
)

type TimeBasedScheduler interface {
	Schedule(time time.Time, filename string)
}

type TimeBasedSchedulerImpl struct {
	Executor exector.Executor
}

func (t *TimeBasedSchedulerImpl) Schedule(timeStr time.Time, filename string) {
	duration := timeStr.Sub(time.Now())
	timer := time.NewTimer(duration)
	done := make(chan bool)
	go func() {
		c := <-timer.C
		log.Printf("Executed %s at %s\n", filename, c.String())
		_, e := t.Executor.ExecuteJob(filename)
		if e != nil {
			log.Fatalf("Unable to execute %s due to %v\n", filename, e)
		}
		done <- true
	}()
	<-done
}
