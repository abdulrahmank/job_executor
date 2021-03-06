package scheduler

import (
	"bytes"
	"fmt"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/internal/mocks"
	"github.com/golang/mock/gomock"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestJobScheduler(t *testing.T) {
	t.Run("Should execute job at a given time", func(t *testing.T) {
		log.Println(time.Now().String())
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockExecutor := mocks.NewMockExecutor(mockCtrl)
		mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
		scheduler := TimeBasedSchedulerImpl{Executor: mockExecutor, SettingDao: mockSettingDao}
		now := time.Now().Add(time.Duration(2 * time.Second))

		mockExecutor.EXPECT().ExecuteJob("hw", "hw.sh")
		mockSettingDao.EXPECT().SetJobStatus("hw", dao.STATUS_SCHEDULED)

		var str bytes.Buffer
		log.SetOutput(&str)

		scheduler.Schedule(now, "hw", "hw.sh")

		expectedLog := fmt.Sprintf("Executed hw.sh at %s\n", now.Format(time.Stamp))
		actualLog := str.String()[strings.Index(str.String(), "Executed"):]

		if actualLog != expectedLog {
			t.Errorf("Expected %s but was %s", expectedLog, actualLog)
		}
		log.SetOutput(os.Stdout)
	})

	t.Run("Should be able to cancel execution of job", func(t *testing.T) {
		log.Println(time.Now().String())
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockExecutor := mocks.NewMockExecutor(mockCtrl)
		mockSettingDao := mocks.NewMockJobSettingDao(mockCtrl)
		scheduler := TimeBasedSchedulerImpl{Executor: mockExecutor, SettingDao: mockSettingDao}
		now := time.Now().Add(time.Duration(2 * time.Second))

		mockExecutor.EXPECT().ExecuteJob(gomock.Any(), gomock.Any()).Times(0)
		mockSettingDao.EXPECT().SetJobStatus(gomock.Any(), gomock.Any())

		go scheduler.Schedule(now, "hw", "hw.sh")
		time.Sleep(time.Second)
		ChannelPool["hw"].CancelCh <- true
	})
}
