package dao

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"testing"
)

func TestJobSettingDaoImpl_SaveJob(t *testing.T) {
	t.Run("Should save given job settings", func(t *testing.T) {
		dao := JobSettingDaoImpl{}
		jobName := "helloWorld"
		daysInWeek := "mon,wed,thu"
		timeSlots := "20:00,10:00"
		fileName := "helloworld.sh"
		numberOfWeeks := 2

		dao.SaveJob(jobName, timeSlots, daysInWeek, fileName, numberOfWeeks)

		db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
		rows, _ := db.Query("SELECT * FROM job_settings WHERE job_name = 'helloWorld'")

		settings := JobSettings{}
		rows.Next()
		if e := rows.Scan(&settings.Id, &settings.JobName, pq.Array(&settings.TimeSlots), pq.Array(&settings.DaysInWeek), &settings.FileName, &settings.NumberOfWeeks); e != nil {
			log.Fatal(e)
		}

		if settings.TimeSlots[0] != "20:00" &&  settings.TimeSlots[1] != "10:00" {
			t.Error("Time mismatch")
			t.Fail()
		}

		if settings.DaysInWeek[0] != "mon" && settings.DaysInWeek[1] != "wed" && settings.DaysInWeek[2] != "thu" {
			t.Error("Days mismatch")
			t.Fail()
		}

		if settings.FileName != fileName {
			t.Error("Filename mismatch")
			t.Fail()
		}

		if settings.NumberOfWeeks != numberOfWeeks {
			t.Error("Number of weeks mismatch")
			t.Fail()
		}

		if settings.JobName != jobName {
			t.Error("Job name mismatch")
			t.Fail()
		}
	})
}