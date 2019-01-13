package dao

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type JobSettingDao interface {
	SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int)
	GetJobFor(day string) []JobSettings
	SetJobStatus(jobName, status string)
}

type JobSettingDaoImpl struct{}

type JobSettings struct {
	Id            int
	JobName       string
	TimeSlots     []string
	DaysInWeek    []string
	FileName      string
	NumberOfWeeks int
}

var db *sql.DB

const STATUS_NOT_PICKED = "not_picked"
const STATUS_SCHEDULED = "scheduled"
const STATUS_COMPLETED = "completed"

func (j *JobSettingDaoImpl) SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int) {
	if e := getDB(); e != nil {
		return
	}
	tx, _ := db.Begin()
	defer tx.Commit()
	timeSlotSlice := strings.Split(timeSlots, ",")
	daysInWeekSlice := strings.Split(daysInWeek, ",")
	result, e := db.Exec(
		"INSERT INTO job_settings (job_name, time_slots, days, file_name, remaining_weeks) VALUES ($1, $2, $3, $4, $5)", jobName, pq.Array(timeSlotSlice), pq.Array(daysInWeekSlice), fileName, numberOfWeeks)
	if e != nil {
		log.Panicf("%v\n", e)
	}
	result, e = db.Exec("INSERT INTO job_status VALUES ($1, $2)", jobName, STATUS_NOT_PICKED)
	if e != nil {
		log.Panicf("%v\n", e)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		log.Fatal("Error inserting job")
	}
}

func (j *JobSettingDaoImpl) GetJobFor(day string) []JobSettings {
	if e := getDB(); e != nil {
		return nil
	}
	rows, e := db.Query("SELECT j.id, j.job_name, j.time_slots, j.days, j.file_name, j.remaining_weeks "+
		"FROM job_settings j LEFT JOIN job_status s ON j.job_name=s.job_name  WHERE days @> ARRAY[$1]::text[] "+
		"AND j.remaining_weeks > 0 AND s.status = $2", day, STATUS_NOT_PICKED)
	if e != nil {
		log.Fatalf("Unable to query for day: %s\n", e.Error())
	}
	jobs := make([]JobSettings, 0)
	for rows.Next() {
		j := JobSettings{}
		if e = rows.Scan(&j.Id, &j.JobName, pq.Array(&j.TimeSlots), pq.Array(&j.DaysInWeek), &j.FileName, &j.NumberOfWeeks); e != nil {
			log.Fatalf("%v\n", e)
		}
		jobs = append(jobs, j)
	}
	return jobs
}

func (j *JobSettingDaoImpl) SetJobStatus(jobName, status string) {
	if e := getDB(); e != nil {
		return
	}
	_, e := db.Exec("UPDATE job_status SET status=$1 WHERE job_name=$2", status, jobName)
	if e != nil {
		log.Panicf("%v\n", e)
	}
}

func getDB() error {
	var e error
	if db == nil {
		db, e = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
		if e != nil {
			log.Panicf("%v\n", e)
		}
	}
	return e
}

func getPSQlInfo(user, dbName, password string) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable password=%s",
		"localhost", 5432, user, dbName, password)
	return psqlInfo
}
