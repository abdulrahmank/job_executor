package orchestrator

import (
	"github.com/abdulrahmank/job_executor/constants"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"github.com/abdulrahmank/job_executor/time_based/scheduler"
	"log"
	"time"
)

type JobOrchestrator struct {
	Scheduler   scheduler.Scheduler
	JobExecutor exector.Executor
	SettingsDao dao.JobSettingDao
}

func (j *JobOrchestrator) SyncJobs() {
	today := time.Now().Weekday().String()
	jobsForToday := j.SettingsDao.GetJobsFor(today)

	for _, job := range jobsForToday {
		for _, timeSlot := range job.TimeSlots {
			j.Scheduler.Schedule(getTime(timeSlot), job.JobName(), job.FileName())
		}
	}
}

func (j *JobOrchestrator) ExecuteJobsForEvent(eventName string) {
	jobsForEvent := j.SettingsDao.GetJobsForEvent(eventName)

	for _, job := range jobsForEvent {
		if _, e := j.JobExecutor.ExecuteJob(job.JobName, job.FileName); e != nil {
			log.Panicf("Error running job %s:%v\n", job.JobName, e)
		}
	}
}

func (j *JobOrchestrator) ExecuteJob(jobName string) {
	fileName := j.SettingsDao.GetFileName(jobName)
	if _, e := j.JobExecutor.ExecuteJob(jobName, fileName); e != nil {
		log.Panicf("Error running job %s:%v\n", jobName, e)
	}
}

func (j *JobOrchestrator) ResetJobStatus() {
	j.SettingsDao.ResetJobStatus(dao.STATUS_COMPLETED, dao.STATUS_NOT_PICKED)
}

func getTime(timeSlot string) time.Time {
	parsedTime, _ := time.Parse(constants.TIME_LAYOUT, timeSlot)
	now := time.Now()
	scheduledTime := time.Date(
		now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.UTC)
	return scheduledTime
}
