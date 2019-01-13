package exector

import (
	"github.com/abdulrahmank/job_executor/constants"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/internal/mocks"
	"github.com/golang/mock/gomock"
	"os"
	"testing"
)

func TestShouldScheduleTimerToExecuteJob(t *testing.T) {
	t.Run("Should execute saved file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettingDao := mocks.NewMockJobSettingDao(ctrl)

		executor := BashExecutorImp{SettingDao: mockSettingDao}
		_ = os.Setenv(constants.JOB_FILE_DIR, "/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/")
		e := os.Chmod("/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/hw.sh", 777)
		if e != nil {
			t.Error(e.Error())
		}
		mockSettingDao.EXPECT().SetJobStatus("hw", dao.STATUS_COMPLETED)

		_, err := executor.ExecuteJob("hw", "hw.sh")

		if err != nil {
			t.Errorf("expected error to be nil but was %v\n", err.Error())
		}

		file, _ := os.Open("/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/file.txt")
		content := make([]byte, 100)
		size, _ := file.Read(content)
		str := string(content[:size])

		expected := "hello world\n"
		if expected != str {
			t.Errorf("expected %s but was %s", expected, str)
		}
		if err != nil {
			t.Errorf("expected error to be nil but was %v\n", err.Error())
		}
	})

	_ = os.Remove("/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/file.txt")
}
