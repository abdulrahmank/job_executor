package dao

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"testing"
)

func TestJobSettingDaoImpl_SaveJob(t *testing.T) {
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
	if e := rows.Scan(&settings.JobName, pq.Array(&settings.TimeSlots), pq.Array(&settings.DaysInWeek), &settings.FileName, &settings.NumberOfWeeks); e != nil {
		log.Fatal(e)
	}

	if settings.TimeSlots[0] != "20:00" && settings.TimeSlots[1] != "10:00" {
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

}

func TestJobSettingDaoImpl_GetJobsFor(t *testing.T) {
	db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
	dao := JobSettingDaoImpl{}
	t.Run("Should get jobs for that day", func(t *testing.T) {
		db.Exec("TRUNCATE job_settings")
		db.Exec("TRUNCATE job_status")

		dao.SaveJob("1", "08:00PM,10:00AM", "mon,wed,thu", "1.sh", 2)
		dao.SaveJob("2", "08:00PM,10:00AM", "tue,fri", "2.sh", 2)
		dao.SaveJob("3", "08:00PM,10:00AM", "wed", "3.sh", 2)

		expectedJobNames := []string{"1", "3"}
		expectedFileNames := []string{"1.sh", "3.sh"}

		jobs := dao.GetJobsFor("wed")

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

	t.Run("Should not fetch job if remaining weeks is 0", func(t *testing.T) {
		db.Exec("TRUNCATE job_settings")
		db.Exec("TRUNCATE job_status")

		dao.SaveJob("1", "08:00PM,10:00AM", "wed", "4.sh", 0)

		jobs := dao.GetJobsFor("wed")

		if len(jobs) != 0 {
			t.Errorf("Count mismatch expected %d but was %d", 0, len(jobs))
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

func TestJobSettingDaoImpl_DecrementRemainingWeeks(t *testing.T) {
	db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
	db.Exec("TRUNCATE job_settings")
	db.Exec("TRUNCATE job_status")
	dao := JobSettingDaoImpl{}
	jobName := "1"
	dao.SaveJob(jobName, "08:00PM,10:00AM", "mon,wed,thu", "1.sh", 2)

	dao.DecrementRemainingWeeks(jobName)

	rows, _ := db.Query("SELECT * FROM job_settings WHERE job_name = $1", jobName)
	rows.Next()
	settings := JobSettings{}
	if e := rows.Scan(&settings.JobName, pq.Array(&settings.TimeSlots), pq.Array(&settings.DaysInWeek), &settings.FileName, &settings.NumberOfWeeks); e != nil {
		log.Fatal(e)
	}
	if settings.NumberOfWeeks != 1 {
		t.Errorf("Expected %d but was %d", 1, settings.NumberOfWeeks)
	}
}

func TestJobSettingDaoImpl_DeleteJob(t *testing.T) {
	db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
	db.Exec("TRUNCATE job_settings")
	db.Exec("TRUNCATE job_status")
	dao := JobSettingDaoImpl{}
	jobName := "1"
	dao.SaveJob(jobName, "08:00PM,10:00AM", "mon,wed,thu", "1.sh", 2)

	dao.DeleteJob(jobName)

	rows, _ := db.Query("SELECT COUNT(*) FROM job_settings WHERE job_name = '1'")
	statusRows, _ := db.Query("SELECT COUNT(*) FROM job_status WHERE job_name = '1'")
	rows.Next()
	statusRows.Next()

	var count int
	if e := rows.Scan(&count); e != nil {
		log.Fatal(e)
	}
	if count != 0 {
		t.Errorf("Expected %d but was %d", 0, count)
	}
	if e := statusRows.Scan(&count); e != nil {
		log.Fatal(e)
	}
	if count != 0 {
		t.Errorf("Expected %d but was %d", 0, count)
	}
}

func TestJobSettingDaoImpl_GetFileName(t *testing.T) {
	db, _ = sql.Open("postgres", getPSQlInfo("test", "test", "password"))
	db.Exec("TRUNCATE job_settings")
	db.Exec("TRUNCATE job_status")
	dao := JobSettingDaoImpl{}
	jobName := "1"
	jobFileName := "1.sh"
	dao.SaveJob(jobName, "08:00PM,10:00AM", "mon,wed,thu", jobFileName, 2)

	fileName := dao.GetFileName(jobName)

	if fileName != jobFileName {
		t.Errorf("Expected %s, but was %s\n", jobFileName, fileName)
	}
}
