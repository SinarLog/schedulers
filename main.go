package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SinarLog/schedulers/jobs"
	"github.com/SinarLog/schedulers/pkg"
)

func main() {
	db := pkg.GetPostgres()
	rdis := pkg.GetRedis()
	scheduler := pkg.GetScheduler()

	if err := jobs.RegisterInitJobs(scheduler, db, rdis); err != nil {
		log.Fatalf("unable to start scheduler: %s\n", err)
	}

	scheduler.StartAsync()
	log.Printf("Scheduler has started")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8123", nil)
}

// For cloud run to health check
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", "World")
}
