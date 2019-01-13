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
		jobs := []dao.JobSettings{{JobName: "job1", TimeSlots: []string{"10:00AM"}, FileName: "1.sh"},
			{JobName: "job2", TimeSlots: []string{"10:00PM"}, FileName: "2.sh"}}
		mockSettingDao.EXPECT().GetJobFor(today).Return(jobs)

		time1 := time.Date(
			now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC)
		time2 := time.Date(
			now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)

		mockScheduler.EXPECT().Schedule(time1, "1.sh")
		mockScheduler.EXPECT().Schedule(time2, "2.sh")

		orchestrator.SyncJobs()
	})
}
