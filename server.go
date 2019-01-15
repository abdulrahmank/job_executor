package main

import (
	"github.com/abdulrahmank/job_executor/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/jobs/create/", handler.JobCreationHandler)
	http.HandleFunc("/jobs/config/", handler.JobConfigHandler)
	http.HandleFunc("/jobs/", handler.JobHandler)
	http.HandleFunc("/event/", handler.EventHandler)
	log.Panic(http.ListenAndServe(":8082", nil))
}
