package orchestrator

import (
	"github.com/abdulrahmank/job_executor/constants"
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
	jobsForToday := j.SettingsDao.GetJobsFor(today)

	for _, job := range jobsForToday {
		for _, timeSlot := range job.TimeSlots {
			j.Scheduler.Schedule(getTime(timeSlot), job.JobName, job.FileName)
		}
	}
}

func getTime(timeSlot string) time.Time {
	parsedTime, _ := time.Parse(constants.TIME_LAYOUT, timeSlot)
	now := time.Now()
	scheduledTime := time.Date(
		now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
	return scheduledTime
}
