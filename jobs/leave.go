package jobs

import (
	"fmt"
	"log"

	"github.com/SinarLog/schedulers/entity"
	"github.com/SinarLog/schedulers/utils"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// LeaveRequestCloserJob closes leave request where the requested
// leave starting date is 3 days after the current time the job
// executes. Any closed leave requests will be marked as closed
// and is not approved by any means. It also generates an automatic
// rejection reason.
func LeaveRequestCloserJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		var leaves []entity.Leave
		var rejectionReason string = "This leave request was closed because it had not finished processing 3 days before the start of request's date."

		// Collect all leaves
		if err := db.WithContext(job.Context()).Model(&entity.Leave{}).
			Where(`
			CAST("leaves"."from" AS date) - CAST(now() AS date) <= 3
			AND
			"leaves"."type" <> 'SICK'
			AND
			parent_id IS NULL
			AND
			closed_automatically IS NULL
			AND
			(
				(
					"leaves"."approved_by_manager" IS NULL
					AND
					"leaves"."approved_by_hr" IS NULL
				)
				OR
				(
					"leaves"."approved_by_manager" IS TRUE
					AND
					"leaves"."approved_by_hr" IS NULL
				)
			)
		`).
			Preload("Childs").
			Find(&leaves).Error; err != nil {
			log.Printf("unable to get leaves for LeaveRequestCloserJob: %s\n", err)
			return
		}

		for _, v := range leaves {
			// Close parent
			tx := db.WithContext(job.Context()).Begin()
			if err := tx.Exec(`UPDATE "leaves" SET closed_automatically = ?, rejection_reason = ? WHERE id = ?`, true, rejectionReason, v.Id).Error; err != nil {
				log.Printf("error while closing parent leave for LeaveRequestCloserJob: %s\n", err)
				tx.Rollback()
				return
			}
			sql := generateUpdateLeaveQuotaSql(v.Type, true)
			if err := tx.Exec(sql, utils.CountNumberOfWorkingDays(v.From, v.To), v.EmployeeID).Error; err != nil {
				log.Printf("error while updating leave quota for LeaveRequestCloserJob: %s\n", err)
				tx.Rollback()
				return
			}

			// Close childs only if pending
			for _, c := range v.Childs {
				if leaveStatusMapper(c) == "PENDING" {
					if err := tx.Exec(`UPDATE "leaves" SET closed_automatically = ?, rejection_reason = ? WHERE id = ?`, true, rejectionReason, c.Id).Error; err != nil {
						log.Printf("error while closing child leave for LeaveRequestCloserJob: %s\n", err)
						tx.Rollback()
						return
					}

					sql := generateUpdateLeaveQuotaSql(c.Type, true)
					if err := tx.Exec(sql, utils.CountNumberOfWorkingDays(c.From, c.To), v.EmployeeID).Error; err != nil {
						log.Printf("error while updating leave quota for LeaveRequestCloserJob: %s\n", err)
						tx.Rollback()
						return
					}
				}
			}

			tx.Commit()
		}

		log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
	}

	_, err := scheduler.Every(1).Day().At("00:00").Tag("leaveCloser").DoWithJobDetails(task)

	return err
}

func OnLeaveStatusJob(scheduler *gocron.Scheduler, db *gorm.DB, redis *redis.Client) error {
	task := func(job gocron.Job) {
		var leaves []entity.Leave

		if err := db.Model(&entity.Leave{}).
			Preload("Employee").
			Where(`"leaves"."from" BETWEEN ? AND ?`,
				utils.GetStartOfDay(),
				utils.GetEndOfDay()).
			Where("approved_by_manager IS TRUE").
			Where("approved_by_hr IS TRUE").
			Find(&leaves).Error; err != nil {
			log.Printf("Unable to query leaves for OnLeaveStatusJob: %s\n", err)
			return
		}

		for _, v := range leaves {
			if v.Employee.Status == entity.RESIGNED {
				continue
			}
			if err := db.WithContext(job.Context()).Exec(`UPDATE "employees" SET status = ? WHERE "employees"."id" = ?`, entity.ON_LEAVE, v.EmployeeID).Error; err != nil {
				log.Printf("unable to update on leave status for OnLeaveStatusJob: %s\n", err)
			}
		}

		log.Printf(SUCCESS_MESSAGE, job.Tags()[0])
	}

	_, err := scheduler.Every(1).Day().At("00:00").Tag("onLeaveStatusUpdater").WaitForSchedule().DoWithJobDetails(task)

	return err
}

func generateUpdateLeaveQuotaSql(leaveType entity.LeaveType, reverse bool) string {
	mapper := map[entity.LeaveType]string{
		entity.ANNUAL:   "yearly_count",
		entity.MARRIAGE: "marriage_count",
		entity.UNPAID:   "unpaid_count",
	}

	switch leaveType {
	case entity.UNPAID:
		if reverse {
			return fmt.Sprintf("UPDATE employee_leaves_quota SET %s = %s - ? WHERE employee_id = ?", mapper[leaveType], mapper[leaveType])
		}
		return fmt.Sprintf("UPDATE employee_leaves_quota SET %s = %s + ? WHERE employee_id = ?", mapper[leaveType], mapper[leaveType])
	case entity.ANNUAL, entity.MARRIAGE:
		if reverse {
			return fmt.Sprintf("UPDATE employee_leaves_quota SET %s = %s + ? WHERE employee_id = ?", mapper[leaveType], mapper[leaveType])
		}
		return fmt.Sprintf("UPDATE employee_leaves_quota SET %s = %s - ? WHERE employee_id = ?", mapper[leaveType], mapper[leaveType])
	default:
		return ""
	}
}

func leaveStatusMapper(v entity.Leave) string {
	var l string

	if v.ApprovedByHr != nil && v.ApprovedByManager != nil {
		if *v.ApprovedByManager && *v.ApprovedByHr {
			l = "APPROVED"
		} else {
			l = "REJECTED"
		}
	} else if v.ApprovedByManager != nil {
		if !*v.ApprovedByManager {
			l = "REJECTED"
		} else {
			l = "PENDING"
		}
	} else {
		l = "PENDING"
	}

	if v.ClosedAutomatically != nil {
		l = "CLOSED"
	}

	return l
}
