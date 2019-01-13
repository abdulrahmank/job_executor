package scheduler

import (
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestJobScheduler(t *testing.T) {
	t.Run("Should execute job at a given time", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockExecutor := mocks.NewMockExecutor(mockCtrl)
		mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
		scheduler := TimeBasedSchedulerImpl{Executor: mockExecutor, SettingDao: mockSettingDao}
		now := time.Now().Add(time.Duration(2 * time.Second))

		mockExecutor.EXPECT().ExecuteJob("hw.sh")
		mockSettingDao.EXPECT().SetJobStatus("hw", dao.STATUS_SCHEDULED)

		scheduler.Schedule(now, "hw", "hw.sh")
	})
}
