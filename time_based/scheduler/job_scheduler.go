package scheduler

import (
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"log"
	"time"
)

type TimeBasedScheduler interface {
	Schedule(time time.Time, jobName, filename string)
}

type SchedulerChannel struct {
	timerCh  <-chan time.Time
	CancelCh chan bool
}
type TimeBasedSchedulerImpl struct {
	SettingDao dao.JobSettingDao
	Executor   exector.Executor
}

var ChannelPool map[string]*SchedulerChannel

func init() {
	ChannelPool = make(map[string]*SchedulerChannel, 5)
}

func (t *TimeBasedSchedulerImpl) Schedule(timeStr time.Time, jobName, filename string) {
	duration := timeStr.Sub(time.Now())
	timer := time.NewTimer(duration)
	channel := SchedulerChannel{timerCh: timer.C, CancelCh: make(chan bool)}
	ChannelPool[jobName] = &channel
	done := make(chan bool)
	go func() {
		select {
		case c := <-channel.timerCh:
			log.Printf("Executed %s at %s\n", filename, c.Format(time.Stamp))
			_, e := t.Executor.ExecuteJob(jobName, filename)
			if e != nil {
				log.Fatalf("Unable to execute %s due to %v\n", filename, e)
			}
			done <- true
			return
		case <-channel.CancelCh:
			done <- true
			return
		}
	}()
	t.SettingDao.SetJobStatus(jobName, dao.STATUS_SCHEDULED)
	<-done
}
