package entity

import (
	"time"
)

type ContractType string

const (
	FULL_TIME ContractType = "FULL_TIME"
	CONTRACT  ContractType = "CONTRACT"
	INTERN    ContractType = "INTERN"
)

type Gender string

const (
	M Gender = "M"
	F Gender = "F"
)

type Religion string

const (
	CHRISTIAN Religion = "CHRISTIAN"
	MUSLIM    Religion = "MUSLIM"
	CATHOLIC  Religion = "CATHOLIC"
	BUDDHA    Religion = "BUDDHA"
	HINDU     Religion = "HINDU"
	CONFUCION Religion = "CONFUCION"
)

type Relation string

const (
	FATHER  Relation = "FATHER"
	MOTHER  Relation = "MOTHER"
	SIBLING Relation = "SIBLING"
	SPOUSE  Relation = "SPOUSE"
)

type Status string

const (
	AVAILABLE   Status = "AVAILABLE"
	UNAVAILABLE Status = "UNAVAILABLE"
	ON_LEAVE    Status = "ON_LEAVE"
	RESIGNED    Status = "RESIGNED"
)

type Employee struct {
	BaseModelId

	FullName     string       `gorm:"type:varchar(155)"`
	Email        string       `gorm:"type:varchar(150);index:,unique,type:btree"`
	Password     string       `gorm:"type:varchar(255)"`
	ContractType ContractType `gorm:"type:varchar(100)"`
	Avatar       string       `gorm:"type:varchar(255)"`
	Status       Status       `gorm:"type:varchar(100);default:'UNAVAILABLE'"`
	IsNewUser    bool
	JoinDate     time.Time
	ResignDate   *time.Time `gorm:"default:null"`

	EmployeeBiodata            EmployeeBiodata
	EmployeesEmergencyContacts []EmployeesEmergencyContact
	EmployeeLeavesQuota        EmployeeLeavesQuota
	EmployeeDataHistoryLogs    []EmployeeDataHistoryLog

	ManagerID    *string `gorm:"type:uuid"`
	Manager      *Employee
	CreatedById  *string `gorm:"type:uuid;default:null"`
	CreatedBy    *Employee
	ResignedById *string `gorm:"type:uuid;default:null"`
	ResignedBy   *Employee
	ResignedAt   *time.Time
	RoleID       string `gorm:"type:uuid"`
	Role         Role
	JobID        string `gorm:"type:uuid"`
	Job          Job

	BaseModelStamps
	BaseModelSoftDelete
}

type EmployeeBiodata struct {
	BaseModelId

	EmployeeID    string   `gorm:"type:uuid,uniqueIndex"`
	NIK           string   `gorm:"type:varchar(255);index:,unique,type:btree"`
	NPWP          string   `gorm:"type:varchar(255);index:,unique,type:btree"`
	Gender        Gender   `gorm:"type:varchar(10)"`
	Religion      Religion `gorm:"type:varchar(85)"`
	PhoneNumber   string   `gorm:"type:varchar(150);index:,unique,type:btree"`
	Address       string
	BirthDate     time.Time
	MaritalStatus bool

	BaseModelStamps
	BaseModelSoftDelete
}

type EmployeesEmergencyContact struct {
	BaseModelId

	EmployeeID  string `gorm:"type:uuid"`
	Employee    Employee
	FullName    string   `gorm:"type:varchar(255)"`
	Relation    Relation `goem:"type:varchar(150)"`
	PhoneNumber string   `gorm:"type:varchar(150)"`

	BaseModelStamps
	BaseModelSoftDelete
}

type EmployeeLeavesQuota struct {
	BaseModelId

	EmployeeID    string `gorm:"type:uuid"`
	YearlyCount   int    `gorm:"default:0"`
	UnpaidCount   int    `gorm:"default:0"`
	MarriageCount int    `gorm:"default:0"`

	BaseModelStamps
	BaseModelSoftDelete
}

type EmployeeDataHistoryLog struct {
	BaseModelId

	EmployeeID  string `gorm:"type:uuid"`
	Employee    Employee
	UpdatedByID string `gorm:"type:uuid"`
	UpdatedBy   Employee
	Changes     JSONB

	BaseModelStamps
	BaseModelSoftDelete
}
