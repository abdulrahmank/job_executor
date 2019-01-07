package dao

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type JobSettingDao interface {
	SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int)
}

type JobSettingDaoImpl struct{}

type JobSettings struct {
	Id         int
	JobName    string
	TimeSlots  []string
	DaysInWeek []string
	FileName   string
	NumberOfWeeks int
}

var db *sql.DB

func (j *JobSettingDaoImpl) SaveJob(jobName, timeSlots, daysInWeek, fileName string, numberOfWeeks int) {
	var e error
	if db == nil {
		db, e = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
		if e != nil {
			log.Panicf("%v\n", e)
		}
	}
	tx, _ := db.Begin()
	defer tx.Commit()
	result, e := tx.Exec(
		"INSERT INTO job_settings (job_name, time_slots, days, file_name, remaining_weeks) VALUES ($1, $2, $3, $4, $5)", jobName, strings.Split(timeSlots, ","), strings.Split(daysInWeek, ","), fileName, numberOfWeeks)
	if e != nil {
		log.Panicf("%v\n", e)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		log.Fatal("Error inserting job")
	}
}

func getPSQlInfo(user, dbName, password string) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable password=%s",
		"localhost", 5432, user, dbName, password)
	return psqlInfo
}

