package dao

import (
	"github.com/abdulrahmank/job_executor/constants"
	"log"
	"os"
)

type FileDao interface {
	SaveFile(jobFileName string, contents []byte)
	DeleteFile(jobFileName string)
}

type FileDaoImpl struct{}

func (f *FileDaoImpl) SaveFile(jobFileName string, contents []byte) {
	jobFilePath := os.Getenv(constants.JOB_FILE_DIR)
	file, e := os.Create(jobFilePath + "/" + jobFileName)
	if e = os.Chmod(jobFileName, 777); e != nil {
		log.Fatal("Unable to change file mode")
	}
	if e != nil {
		log.Fatal("Unable to create file")
	}
	if _, e = file.Write(contents); e != nil {
		log.Fatal("Unable to write file contents")
	}
	if e = file.Close(); e != nil {
		log.Fatal("Unable to save file contents")
	}
}

func (f *FileDaoImpl) DeleteFile(jobFileName string) {
	jobFilePath := os.Getenv(constants.JOB_FILE_DIR)
	e := os.Remove(jobFilePath + "/" + jobFileName)
	if e != nil {
		log.Fatal("Unable to create file")
	}
}
