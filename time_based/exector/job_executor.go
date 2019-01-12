package exector

import (
	"github.com/abdulrahmank/job_executor/constants"
	"log"
	"os"
	"os/exec"
)

type Executor interface {
	ExecuteJob(jobName, fileName string) ([]byte, error)
}

type BashExecutorImp struct{}

func (e *BashExecutorImp) ExecuteJob(fileName string) ([]byte, error) {
	jobFileDir := os.Getenv(constants.JOB_FILE_DIR)
	cmd := exec.Command("/bin/sh", jobFileDir+fileName)
	if op, err := cmd.Output(); err != nil {
		log.Fatal(err.Error())
		return op, err
	}
	return nil, nil
}
