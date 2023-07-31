package jobs

import (
	"log"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func HealthCheck(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		// Check db conn
		if err := db.WithContext(job.Context()).Raw("SELECT 1=1 AS t").Error; err != nil {
			log.Printf("db conn problem for HealthCheck: %s\n", err)
			return
		}

		// Check redis conn
		if err := redis.Ping(job.Context()).Err(); err != nil {
			log.Printf("redis conn problem for HealthCheck: %s\n", err)
			return
		}

		log.Printf("health check completed")
	}

	_, err := scheduler.Every(30).Minute().Tag("healthCheck").DoWithJobDetails(task)

	return err
}
