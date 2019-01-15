package orchestrator

import (
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestJobOrchestrator(t *testing.T) {
	t.Run("Should fetch jobs and schedule them", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
		mockScheduler := mocks.NewMockTimeBasedScheduler(mockCtrl)
		orchestrator := JobOrchestrator{Scheduler: mockScheduler,
			SettingsDao: mockSettingDao}
		now := time.Now()
		today := now.Weekday().String()
		jobs := []dao.TimeBasedJob{{JobVal: &dao.Job{JobName: "job1", FileName: "1.sh"}, TimeSlots: []string{"10:00AM"}},
			{JobVal: &dao.Job{JobName: "job2", FileName: "2.sh"}, TimeSlots: []string{"10:00PM"}}}
		mockSettingDao.EXPECT().GetJobsFor(today).Return(jobs)

		time1 := time.Date(
			now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC)
		time2 := time.Date(
			now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)

		mockScheduler.EXPECT().Schedule(time1, "job1", "1.sh")
		mockScheduler.EXPECT().Schedule(time2, "job2", "2.sh")

		orchestrator.SyncJobs()
	})
}

func TestJobOrchestrator_ExecuteJobsForEvent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockExecutor := mocks.NewMockExecutor(mockCtrl)
	mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
	orchestrator := JobOrchestrator{JobExecutor: mockExecutor, SettingsDao: mockSettingDao}
	eventName := "ping"
	jobName := "1"
	fileName := "1.sh"
	jobs := []dao.Job{{JobName: jobName, FileName: fileName}}

	mockSettingDao.EXPECT().GetJobsForEvent(eventName).Return(jobs)
	mockExecutor.EXPECT().ExecuteJob(jobName, fileName)

	orchestrator.ExecuteJobsForEvent(eventName)
}

func TestJobOrchestrator_ExecuteJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockExecutor := mocks.NewMockExecutor(mockCtrl)
	mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
	jobName := "job"
	fileName := "fileName"

	orchestrator := JobOrchestrator{JobExecutor: mockExecutor, SettingsDao: mockSettingDao}
	mockSettingDao.EXPECT().GetFileName(jobName).Return(fileName)
	mockExecutor.EXPECT().ExecuteJob(jobName, fileName)

	orchestrator.ExecuteJob(jobName)
}
