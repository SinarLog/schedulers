package entity

import (
	"time"
)

type LeaveType string

const (
	ANNUAL   LeaveType = "ANNUAL"
	UNPAID   LeaveType = "UNPAID"
	SICK     LeaveType = "SICK"
	MARRIAGE LeaveType = "MARRIAGE"
)

func (l LeaveType) String() string {
	switch l {
	case ANNUAL:
		return "ANNUAL"
	case UNPAID:
		return "UNPAID"
	case SICK:
		return "SICK"
	case MARRIAGE:
		return "MARRIAGE"
	default:
		return ""
	}
}

type Leave struct {
	BaseModelId

	EmployeeID    string `gorm:"type:uuid"`
	Employee      Employee
	From          time.Time
	To            time.Time
	Type          LeaveType `gorm:"type:varchar(100)"`
	Reason        string    `gorm:"type:text"`
	AttachmentUrl string    `gorm:"type:varchar(255)"`

	// A parent leave contains the original leave request.
	Parent   *Leave
	ParentID *string `gorm:"type:uuid;default:null"`
	// A childs leave contains the overflowed excess of the
	// original leave request. For example, requesting a leave
	// of type ANNUAL with 14 days duration meanwhile I have only
	// 12 ANNUAL quota left. Hence, I can overflow it to an
	// UNPAID leave request of 2 days, with my parent leave request
	// of ANNUAL of 12 days.
	Childs []Leave `gorm:"foreignKey:ParentID"`

	ManagerID           *string `gorm:"type:uuid;default:null"`
	Manager             *Employee
	HrID                *string `gorm:"type:uuid;default:null"`
	Hr                  *Employee
	ApprovedByManager   *bool
	ApprovedByHr        *bool
	ActionByManagerAt   *time.Time
	ActionByHrAt        *time.Time
	RejectionReason     string `gorm:"type:text"`
	ClosedAutomatically *bool

	BaseModelStamps
	BaseModelSoftDelete
}

type LeaveReport struct {
	// Whether the leaves quota exceeds the quota according to the type
	IsLeaveLeakage bool
	// The excessive leave duration as days
	ExcessLeaveDuration int

	// Stores the initial request type
	RequestType LeaveType
	// Stores the remaining quota for the request type
	RemainingQuotaForRequestedType int

	// The available excess types to overflow the leakage
	AvailableExcessTypes []LeaveType
	// The available excess quota to overflow the leakage
	// NOTES: Make unpaid count to 10 max
	AvailableExcessQuotas []int
}
