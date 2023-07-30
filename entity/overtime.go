package entity

import (
	"time"
)

type Overtime struct {
	BaseModelId

	AttendanceID  string `gorm:"type:uuid"`
	Duration      int
	Reason        string `gorm:"type:text"`
	AttachmentUrl string `gorm:"type:varchar(255)"`

	ManagerID           *string
	Manager             *Employee
	ApprovedByManager   *bool
	ActionByManagerAt   *time.Time
	RejectionReason     string
	ClosedAutomatically *bool

	Attendance Attendance

	BaseModelStamps
	BaseModelSoftDelete
}

type OvertimeOnAttendanceReport struct {
	// Whether an attendance is an overtime
	IsOvertime bool `json:"isOvertime"`
	// Whether the attendance made is on holiday
	IsOnHoliday bool `json:"isOnHoliday"`
	// Whether the attendance duration is more than the allowed daily/weekly overtime duration
	IsOvertimeLeakage bool `json:"isOvertimeLeakage"`
	// Whether there could be made an overtime for that week
	IsOvertimeAvailable bool `json:"isOvertimeAvailable"`
	// Attendance's overtime duration
	OvertimeDuration time.Duration `json:"overtimeDuration"`
	// Overtime total duration for this week
	OvertimeWeekTotalDuration time.Duration `json:"overtimeWeeklyTotalDuration"`
	// Overtime accepted duration
	OvertimeAcceptedDuration time.Duration `json:"overtimeAcceptedDuration"`
	// Max allowed overtime daily duration
	MaxAllowedDailyDuration time.Duration `json:"maxAllowedDailyDuration,omitempty"`
	// Max allowed overtime weekly duration
	MaxAllowedWeeklyDuration time.Duration `json:"maxAllowedWeeklyDuration,omitempty"`
}
