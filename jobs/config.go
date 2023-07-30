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

func UpdateConfigNextDayJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		// Checks if the key exists
		exists, err := redis.Exists(job.Context(), entity.CONFIG_KEY_NEXT_DAY).Result()
		if err != nil {
			log.Printf("Unable to check if key exists in redis for *schedulerService.(UpdateConfigNextDayJob): %s\n", err)
			return
		}

		if exists == 1 {
			var configToChange entity.Configuration
			var config entity.Configuration

			// Get config to change from redis
			if err := redis.HGetAll(job.Context(), entity.CONFIG_KEY_NEXT_DAY).Scan(&configToChange); err != nil {
				log.Printf("Unable to retrieve config hash from redis: %s\n", err)
				return
			}

			// Query current config
			tx := db.WithContext(job.Context()).Begin()
			if err := tx.Model(&config).First(&config).Error; err != nil {
				log.Printf("Unable to get config: %s\n", err)
				tx.Rollback()
				return
			}

			// Apply changes
			config.AcceptanceAttendanceInterval = configToChange.AcceptanceAttendanceInterval
			config.OfficeStartTime = time.Date(
				config.OfficeStartTime.Year(),
				config.OfficeStartTime.Month(),
				config.OfficeStartTime.Day(),
				configToChange.OfficeStartTimeHour,
				configToChange.OfficeStartTimeMinute,
				0,
				0,
				utils.CURRENT_LOC,
			)
			config.OfficeEndTime = time.Date(
				config.OfficeEndTime.Year(),
				config.OfficeEndTime.Month(),
				config.OfficeEndTime.Day(),
				configToChange.OfficeEndTimeHour,
				configToChange.OfficeEndTimeMinute,
				0,
				0,
				utils.CURRENT_LOC,
			)

			// Save changes
			if err := tx.Model(&config).Save(&config).Error; err != nil {
				log.Printf("unable to save config changes for UpdateConfigNextDayJob: %s\n", err)
				tx.Rollback()
				return
			}

			if redis.Del(job.Context(), entity.CONFIG_KEY_NEXT_DAY).Val() == 0 {
				log.Printf("unable to delete redis key for UpdateConfigNextDayJob: %s\n", err)
			}
			tx.Commit()
			log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
			return
		}
	}

	_, err := scheduler.Every(1).Day().At("00:00").Tag("configNextDayUpdater").WaitForSchedule().DoWithJobDetails(task)

	return err
}

func UpdateConfigNextMonthJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		// Checks if the key exists
		exists, err := redis.Exists(job.Context(), entity.CONFIG_KEY_NEXT_MONTH).Result()
		if err != nil {
			log.Printf("unable to check if key exists UpdateConfigNextMonthJob: %s\n", err)
			return
		}

		if exists == 1 {
			var configToChange entity.Configuration
			var config entity.Configuration

			// Get config to change from redis
			if err := redis.HGetAll(job.Context(), entity.CONFIG_KEY_NEXT_MONTH).Scan(&configToChange); err != nil {
				log.Printf("unable to get config hash for UpdateConfigNextMonthJob: %s\n", err)
				return
			}

			// Query current config
			tx := db.WithContext(job.Context()).Begin()
			if err := tx.Model(&config).First(&config).Error; err != nil {
				log.Printf("unable to get config record for UpdateConfigNextMonthJob: %s\n", err)
				tx.Rollback()
				return
			}

			// Apply changes
			config.AcceptanceLeaveInterval = configToChange.AcceptanceLeaveInterval
			config.DefaultYearlyQuota = configToChange.DefaultYearlyQuota
			config.DefaultMarriageQuota = configToChange.DefaultMarriageQuota

			// Save changes
			if err := tx.Model(&config).Save(&config).Error; err != nil {
				log.Printf("unable to save config changes for UpdateConfigNextMonthJob: %s\n", err)
				tx.Rollback()
				return
			}

			if redis.Del(job.Context(), entity.CONFIG_KEY_NEXT_MONTH).Val() == 0 {
				log.Printf("unable to delete redis key for UpdateConfigNextMonthJob %s\n", err)
			}

			tx.Commit()

			log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
		}
	}

	_, err := scheduler.Cron("0 0 1 * *").Tag("configNextMonthUpdater").WaitForSchedule().DoWithJobDetails(task)

	return err
}
