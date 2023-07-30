package entity

import (
	"time"
)

type Attendance struct {
	BaseModelId

	EmployeeID          string `gorm:"type:uuid"`
	Employee            Employee
	ClockInAt           time.Time `gorm:"default:now()"`
	ClockOutAt          time.Time
	DoneForTheDay       bool
	ClockInLoc          Point
	ClockOutLoc         Point
	LateClockIn         bool
	EarlyClockOut       bool
	ClosedAutomatically *bool

	Overtime *Overtime

	BaseModelStamps
	BaseModelSoftDelete
}
