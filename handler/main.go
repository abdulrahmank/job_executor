package handler

import (
	"github.com/abdulrahmank/job_executor/time_based"
	"log"
	"net/http"
	"strconv"
)

func TimedJobScheduler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		if err := r.ParseForm(); err != nil {
			log.Println("Couldn't parse form data")
		}
		jobName := r.FormValue("name")
		timeSlots := r.FormValue("time")
		daysInWeek := r.FormValue("days")
		numberOfWeeks, err := strconv.Atoi(r.FormValue("weeks"))
		if err != nil {
			log.Println("Invalid number of weeks")
		}
		file, _ , err := r.FormFile("script")
		if err != nil {
			log.Println("Error parsing file")
		}
		time_based.SaveJob(jobName, timeSlots, daysInWeek, numberOfWeeks, file)
		break
	case "GET":
		log.Println("Not implemented")
		break
	case "DELETE":
		log.Println("Not implemented")
		break
	}


}