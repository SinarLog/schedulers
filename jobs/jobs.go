package jobs

import (
	"github.com/SinarLog/schedulers/utils"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const SUCCESS_MESSAGE = "Successfully ran job %s\n"

func RegisterInitJobs(scheduler *gocron.Scheduler, db *gorm.DB, rdis *redis.Client) error {
	var errs error
	// Register init jobs
	jobs := []func(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error{
		HealthCheck,
		AttendanceCloserJob,
		LeaveRequestCloserJob,
		OnLeaveStatusJob,
		UpdateConfigNextDayJob,
		UpdateConfigNextMonthJob,
		OvertimeCloserJob,
	}

	for _, f := range jobs {
		if err := f(scheduler, db, rdis); err != nil {
			utils.AddError(errs, err)
		}
	}

	return errs
}
