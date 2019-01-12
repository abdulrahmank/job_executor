package scheduler

import (
	"github.com/abdulrahmank/job_executor/constants"
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"log"
	"time"
)

type TimeBasedScheduler struct {
	Executor exector.Executor
}

func (t *TimeBasedScheduler) Schedule(timeStr, filename string) {
	parse, err := time.Parse(constants.TIME_LAYOUT, timeStr)
	if err != nil {
		log.Fatalf("Unable to parse time %s", err.Error())
	}
	duration := parse.Sub(time.Now())
	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		_, e := t.Executor.ExecuteJob(filename)
		if e != nil {
			log.Fatalf("Unable to execute %s due to %v\n", filename, e)
		}
	}()

}
