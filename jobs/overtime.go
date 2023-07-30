package jobs

import (
	"log"

	"github.com/SinarLog/schedulers/entity"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func OvertimeCloserJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		var overtimes []entity.Overtime
		var truee bool = true

		// Query all the overtimes that has not been processed
		if err := db.WithContext(job.Context()).
			Model(&entity.Overtime{}).
			Where("approved_by_manager IS NULL").
			Where("action_by_manager_at IS NULL").
			Preload("Attendance.Employee").
			Find(&overtimes).Error; err != nil {
			log.Printf("unable to query all pending overtimes for OvertimeJobCloser: %s\n", err)
			return
		}

		for _, v := range overtimes {
			tx := db.WithContext(job.Context()).Begin()

			v.ClosedAutomatically = &truee
			v.RejectionReason = "This overtime submission is closed because it was not processed until 24th of the month."

			if err := tx.Model(&v).Save(&v).Error; err != nil {
				log.Printf("unable to save overtimme for OvertimeJobCloser: %s\n", err)
				tx.Rollback()
				return
			}

			tx.Commit()
		}

		log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
	}

	_, err := scheduler.Cron("0 0 24 * *").Tag("overtimeCloser").WaitForSchedule().DoWithJobDetails(task)

	return err
}
