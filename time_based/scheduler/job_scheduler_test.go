package scheduler

import (
	"github.com/abdulrahmank/job_executor/time_based/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestJobScheduler(t *testing.T) {
	t.Run("Should execute job in given time", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockExecutor := mocks.NewMockExecutor(mockCtrl)
		scheduler := TimeBasedScheduler{Executor: mockExecutor}

		mockExecutor.EXPECT().ExecuteJob("hw.sh")

		scheduler.Schedule("10:00", "hw.sh")
	})
}
