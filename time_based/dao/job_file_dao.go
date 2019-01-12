package dao

import (
	"log"
	"os"
)

type FileDao interface {
	SaveFile(jobFileName string, contents []byte)
}

type FileDaoImpl struct {}

func (f *FileDaoImpl) SaveFile(jobFileName string, contents []byte) {
	file, e := os.Create(jobFileName)
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
