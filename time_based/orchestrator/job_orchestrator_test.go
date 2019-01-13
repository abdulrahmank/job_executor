package orchestrator

import (
	"github.com/abdulrahmank/job_executor/constants"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/mocks"
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
		timeStr := now.Format(constants.TIME_LAYOUT)
		jobs := []dao.JobSettings{{JobName: "job1", TimeSlots: []string{timeStr}, FileName: "1.sh"},
			{JobName: "job2", TimeSlots: []string{timeStr}, FileName: "2.sh"}}
		mockSettingDao.EXPECT().GetJobFor(today).Return(jobs)

		mockScheduler.EXPECT().Schedule(timeStr, "1.sh")
		mockScheduler.EXPECT().Schedule(timeStr, "2.sh")

		orchestrator.SyncJobs()
	})
}
