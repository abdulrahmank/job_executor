package main

import (
	"github.com/abdulrahmank/job_executor/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/timed/", handler.TimedJobScheduler)
	log.Panic(http.ListenAndServe(":8082", nil))
}