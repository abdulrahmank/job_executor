package exector

import (
	"github.com/abdulrahmank/job_executor/constants"
	"log"
	"os"
	"testing"
)

func TestShouldScheduleTimerToExecuteJob(t *testing.T) {
	t.Run("Should execute saved file", func(t *testing.T) {
		executor := BashExecutorImp{}
		_ = os.Setenv(constants.JOB_FILE_DIR, "/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/")
		e := os.Chmod("/Volumes/abdul_extended/go/src/github.com/abdulrahmank/job_executor/test_resource/hw.sh", 777)
		if e != nil {
			t.Error(e.Error())
		}

		op, err := executor.ExecuteJob("hw.sh")
		log.Println(string(op))
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
