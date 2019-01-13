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

		db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
		db.Exec("TRUNCATE job_settings")
		db.Exec("TRUNCATE job_status")

		dao.SaveJob(jobName, timeSlots, daysInWeek, fileName, numberOfWeeks)

		rows, _ := db.Query("SELECT * FROM job_settings WHERE job_name = 'helloWorld'")
		statusRows, _ := db.Query("SELECT * FROM job_status")

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
		}

		statusRows.Next()
		var jobNameStatus, status string
		if e := statusRows.Scan(&jobNameStatus, &status); e != nil {
			log.Fatalf("%v\n", e)
		}
		if jobNameStatus != settings.JobName {
			t.Error("Job name mismatch")
		}
		if status != STATUS_NOT_PICKED {
			t.Error("Job status mismatch")
		}
	})


	t.Run("Should get saved job settings for that day", func(t *testing.T) {
		db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
		db.Exec("TRUNCATE job_settings")
		db.Exec("TRUNCATE job_status")
		dao := JobSettingDaoImpl{}
		dao.SaveJob("1", "08:00PM,10:00AM", "mon,wed,thu", "1.sh", 2)
		dao.SaveJob("2", "08:00PM,10:00AM", "tue,fri", "2.sh", 2)
		dao.SaveJob("3", "08:00PM,10:00AM", "wed", "3.sh", 2)
		dao.SaveJob("4", "08:00PM,10:00AM", "wed", "4.sh", 0)
		expectedJobNames := []string{"1", "3"}
		expectedFileNames := []string{"1.sh", "3.sh"}

		jobs := dao.GetJobFor("wed")

		if len(jobs) != 2 {
			t.Errorf("Count mismatch expected %d but was %d", 2, len(jobs))
		}

		for i, job := range jobs {
			if job.JobName != expectedJobNames[i] {
				t.Errorf("Name mismatch expected %s but was %s", expectedJobNames[i], job.JobName)
			}
			if job.FileName != expectedFileNames[i] {
				t.Errorf("File name mismatch expected %s but was %s", expectedJobNames[i], job.JobName)
			}
		}
	})
}

func TestJobSettingDaoImpl_SetJobStatus(t *testing.T) {
	db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
	db.Exec("TRUNCATE job_settings")
	db.Exec("TRUNCATE job_status")
	dao := JobSettingDaoImpl{}
	dao.SaveJob("1", "08:00PM,10:00AM", "mon,wed,thu", "1.sh", 2)

	dao.SetJobStatus("1", STATUS_SCHEDULED)

	statusRows, _ := db.Query("SELECT * FROM job_status")

	statusRows.Next()
	var jobNameStatus, status string
	if e := statusRows.Scan(&jobNameStatus, &status); e != nil {
		log.Fatalf("%v\n", e)
	}
	if jobNameStatus != "1" {
		t.Error("Job name mismatch")
	}
	if status != STATUS_SCHEDULED {
		t.Errorf("Expected %s but was %s", STATUS_SCHEDULED, status)
	}
}
