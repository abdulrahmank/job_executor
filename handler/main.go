package handler

import (
	"encoding/json"
	"github.com/abdulrahmank/job_executor/time_based/dao"
	"github.com/abdulrahmank/job_executor/time_based/exector"
	"github.com/abdulrahmank/job_executor/time_based/orchestrator"
	"github.com/abdulrahmank/job_executor/time_based/persistor"
	"github.com/abdulrahmank/job_executor/time_based/scheduler"
	"log"
	"net/http"
	"time"
)

type Config struct {
	JobType       *string `json:"type"`
	JobName       *string `json:"jobName"`
	EventName     *string `json:"evenName"`
	TimeSlots     *string `json:"timeSlots"`
	DaysInWeek    *string `json:"days"`
	NumberOfWeeks *int    `json:"week"`
}

type jobActionRequest struct {
	action string `json:"action"`
}

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

func JobCreationHandler(w http.ResponseWriter, r *http.Request) {
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
	}
}

func JobHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		jobName := r.URL.Path[len("/jobs/"):]
		decoder := json.NewDecoder(r.Body)
		actionReq := jobActionRequest{}
		if err := decoder.Decode(&actionReq); err != nil {
			log.Panicf("Error parsing json %v", err)
		}
		switch actionReq.action {
		case "stop":
			scheduler.ChannelPool[jobName].CancelCh <- true
			break
		}
		w.WriteHeader(http.StatusOK)
		break
	case "GET":
		log.Println("Not implemented")
		break
	case "DELETE":
		jobName := r.URL.Path[len("/jobs/"):]
		scheduler.ChannelPool[jobName].CancelCh <- true
		jobPersistor.DeleteJob(jobName)
		w.WriteHeader(http.StatusOK)
		break
	}
}

func JobConfigHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		config := Config{}
		if err := decoder.Decode(&config); err != nil {
			log.Panicf("Error parsing config json %v", err)
		}
		switch *config.JobType {
		case "time":
			jobPersistor := getPersistor()
			if config.JobName == nil || config.TimeSlots == nil || config.DaysInWeek == nil ||
				config.NumberOfWeeks == nil {
				log.Panic("Config options can't be nil")
			}
			jobPersistor.ConfigureTimeBasedJob(*config.JobName, *config.TimeSlots, *config.DaysInWeek,
				*config.NumberOfWeeks)
			jobOrchestrator := getOrchestrator()
			jobOrchestrator.SyncJobs()
			break
		case "event":
			jobPersistor := getPersistor()
			if config.JobName == nil || config.EventName == nil {
				log.Panic("Config options can't be nil")
			}
			jobPersistor.ConfigureEventBasedJob(*config.JobName, *config.EventName)
			break
		}

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
