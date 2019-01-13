package orchestrator

import (
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/scheduler"
	"time"
)

type JobOrchestrator struct {
	Scheduler   scheduler.TimeBasedScheduler
	SettingsDao dao.JobSettingDao
}

func (j *JobOrchestrator) SyncJobs() {
	today := time.Now().Weekday().String()
	jobsForToday := j.SettingsDao.GetJobFor(today)

	for _, job := range jobsForToday {
		for _, timeSlot := range job.TimeSlots {
			j.Scheduler.Schedule(timeSlot, job.FileName)
		}
	}
}
