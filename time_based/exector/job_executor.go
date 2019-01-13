package exector

import (
	"github.com/abdulrahmank/job_executor/constants"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"log"
	"os"
	"os/exec"
)

type Executor interface {
	ExecuteJob(jobName, fileName string) ([]byte, error)
}

type BashExecutorImp struct {
	SettingDao dao.JobSettingDao
}

func (e *BashExecutorImp) ExecuteJob(jobName, fileName string) ([]byte, error) {
	jobFileDir := os.Getenv(constants.JOB_FILE_DIR)
	cmd := exec.Command("/bin/sh", jobFileDir+fileName)
	e.SettingDao.SetJobStatus(jobName, dao.STATUS_COMPLETED)
	if op, err := cmd.Output(); err != nil {
		log.Fatal(err.Error())
		return op, err
	}
	return nil, nil
}
