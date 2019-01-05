package handler

import (
	"log"
	"net/http"
)

func TimedJobScheduler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		log.Println("Not implemented")
		break
	case "GET":
		log.Println("Not implemented")
		break
	case "DELETE":
		log.Println("Not implemented")
		break
	}


}