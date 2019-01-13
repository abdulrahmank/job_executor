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

func (j *JobSettingDaoImpl) SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int) {
	if e := getDB(); e != nil {
		return
	}
	tx, _ := db.Begin()
	defer tx.Commit()
	timeSlotSlice := strings.Split(timeSlots, ",")
	daysInWeekSlice := strings.Split(daysInWeek, ",")
	result, e := tx.Exec(
		"INSERT INTO job_settings (job_name, time_slots, days, file_name, remaining_weeks) VALUES ($1, $2, $3, $4, $5)", jobName, pq.Array(timeSlotSlice), pq.Array(daysInWeekSlice), fileName, numberOfWeeks)
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
	rows, e := db.Query("SELECT * FROM job_settings WHERE days @> ARRAY[$1]::text[] AND remaining_weeks > 0", day)
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
