package handler

import (
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"github.com/abdulrahmank/job_executor/time_based/orchestrator"
	"github.com/abdulrahmank/job_executor/time_based/persistor"
	"github.com/abdulrahmank/job_executor/time_based/scheduler"
	"log"
	"net/http"
	"time"
)

var jobOrchestrator *orchestrator.JobOrchestrator
var jobPersistor *persistor.Persistor

func init() {
	jobOrchestrator = getOrchestrator()
	jobPersistor = getPersistor()

	jobOrchestrator.SyncJobs()
	go syncJobsDaily()
}

func syncJobsDaily() {
	for {
		now := time.Now()
		tmr := now.AddDate(0, 0, 1)
		syncScheduleTime := time.Date(tmr.Year(), tmr.Month(), tmr.Day(),
			0, 0, 0, 0, time.UTC)
		durationForNextSync := syncScheduleTime.Sub(now)
		syncScheduleCh := time.After(durationForNextSync)
		done := make(chan bool)
		select {
		case <-syncScheduleCh:
			getOrchestrator().SyncJobs()
			done <- true
			break
		}
		<-done
	}
}

func JobHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		if err := r.ParseForm(); err != nil {
			log.Println("Couldn't parse form data")
		}
		jobName := r.FormValue("name")
		file, fileHeader, err := r.FormFile("script")
		if err != nil {
			log.Println("Error parsing file")
		}
		var contents []byte
		if _, err = file.Read(contents); err != nil {
			log.Println("Error reading file contents")
		}
		jobPersistor.SaveJob(jobName, fileHeader.Filename, contents)

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Job saved successfully")); err != nil {
			log.Fatal("Unable to wrtie to response")
		}
		break
	case "GET":
		log.Println("Not implemented")
		break
	case "DELETE":
		jobName := r.URL.Path[len("/jobs/"):]
		scheduler.ChannelPool[jobName].CancelCh <- true
		jobPersistor.DeleteJob(jobName)
		break
	}

}

func getOrchestrator() *orchestrator.JobOrchestrator {
	return &orchestrator.JobOrchestrator{Scheduler:
	&scheduler.TimeBasedSchedulerImpl{Executor: &exector.BashExecutorImp{}},
		SettingsDao: &dao.JobSettingDaoImpl{}}

}

func getPersistor() *persistor.Persistor {
	return &persistor.Persistor{FileDao: &dao.FileDaoImpl{}, SettingDao: &dao.JobSettingDaoImpl{}}
}
