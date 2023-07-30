package jobs

import (
	"log"
	"time"

	"github.com/SinarLog/schedulers/entity"
	"github.com/SinarLog/schedulers/utils"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AttendanceCloserJob closes an active attendances made by
// employees which runs every day and executes at 9:00 PM.
// On a closed attendance, it does not create overtime.
func AttendanceCloserJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		now := time.Now().In(utils.CURRENT_LOC)
		truee := true
		tx := db.WithContext(job.Context()).Begin()

		if err := tx.Exec(`UPDATE "attendances" SET clock_out_at = ?, closed_automatically = ?, updated_at = ?, done_for_the_day = ? WHERE done_for_the_day = ? OR clock_out_at = ?`,
			now,
			&truee,
			now,
			true,
			false,
			time.Time{},
		).Error; err != nil {
			log.Printf("unable to exec query for AttendanceCloserJob: %s\n", err)
			tx.Rollback()
			return
		}

		if err := tx.Exec(`UPDATE "employees" SET status = ? WHERE status = ? AND (resigned_at IS NULL OR resigned_by_id IS NULL)`, entity.UNAVAILABLE, entity.AVAILABLE).Error; err != nil {
			log.Printf("unable to exec query for AttendanceCloserJob: %s\n", err)
			tx.Rollback()
			return
		}

		tx.Commit()

		log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
	}

	_, err := scheduler.Every(1).Day().At("21:00").Tag("attendanceCloser").WaitForSchedule().DoWithJobDetails(task)

	return err
}
