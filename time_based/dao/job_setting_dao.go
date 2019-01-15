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
	GetJobsFor(day string) []JobSettings
	SetJobStatus(jobName, status string)
	DecrementRemainingWeeks(jobName string)
	DeleteJob(jobName string)
	GetFileName(jobName string) string
}

type JobSettingDaoImpl struct{}

type JobSettings struct {
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

//TODO: Think of separating remaining weeks from job_settings table
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

func (j *JobSettingDaoImpl) GetJobsFor(day string) []JobSettings {
	if e := getDB(); e != nil {
		return nil
	}
	rows, e := db.Query("SELECT j.job_name, j.time_slots, j.days, j.file_name, j.remaining_weeks "+
		"FROM job_settings j LEFT JOIN job_status s ON j.job_name=s.job_name  WHERE days @> ARRAY[$1]::text[] "+
		"AND j.remaining_weeks > 0 AND s.status = $2", day, STATUS_NOT_PICKED)
	if e != nil {
		log.Fatalf("Unable to query for day: %s\n", e.Error())
	}
	jobs := make([]JobSettings, 0)
	for rows.Next() {
		j := JobSettings{}
		if e = rows.Scan(&j.JobName, pq.Array(&j.TimeSlots), pq.Array(&j.DaysInWeek), &j.FileName, &j.NumberOfWeeks); e != nil {
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

func (j *JobSettingDaoImpl) DecrementRemainingWeeks(jobName string) {
	if e := getDB(); e != nil {
		return
	}
	_, e := db.Exec("update job_settings set remaining_weeks = "+
		"(select remaining_weeks from job_settings where job_name=$1) - 1 where job_name=$1", jobName)
	if e != nil {
		log.Panicf("%v\n", e)
	}
}

func (j *JobSettingDaoImpl) DeleteJob(jobName string) {
	if e := getDB(); e != nil {
		return
	}
	_, e := db.Exec("DELETE FROM job_settings WHERE job_name=$1", jobName)
	if e != nil {
		log.Panicf("%v\n", e)
	}
	_, e = db.Exec("DELETE FROM job_status WHERE job_name=$1", jobName)
	if e != nil {
		log.Panicf("%v\n", e)
	}
}

func (j *JobSettingDaoImpl) GetFileName(jobName string) string {
	if e := getDB(); e != nil {
		return ""
	}
	rows, e := db.Query("SELECT file_name FROM job_settings WHERE job_name=$1", jobName)
	if e != nil {
		log.Fatalf("Unable to query %s\n", e.Error())
	}
	rows.Next()
	var result string
	if e = rows.Scan(&result); e != nil {
		log.Fatalf("Unable to scan %s\n", e.Error())
	}
	return result
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
