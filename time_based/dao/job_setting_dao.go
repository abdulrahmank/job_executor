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

func getPSQlInfo(user, dbName, password string) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable password=%s",
		"localhost", 5432, user, dbName, password)
	return psqlInfo
}
