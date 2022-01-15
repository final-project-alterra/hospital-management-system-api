package schedules

import "time"

type PrescriptionCore struct {
	ID          int
	Medicine    string
	Instruction string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OutpatientCore struct {
	ID        int
	Complaint string
	Status    int
	StartTime string
	EndTime   string
	CreatedAt time.Time
	UpdatedAt time.Time

	WorkSchedule  WorkScheduleCore
	Prescriptions []PrescriptionCore
	Patient       PatientCore
}

type WorkScheduleCore struct {
	ID           int
	Group        string // uuid
	Date         string
	StartTime    string
	EndTime      string
	TotalWaiting int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Nurse        NurseCore
	Doctor       DoctorCore

	Outpatients []OutpatientCore
}

type PatientCore struct {
	ID      int
	NIK     string
	Name    string
	Age     int
	Phone   string
	Address string
	Gender  string
}

type NurseCore struct {
	ID     int
	Email  string
	Name   string
	Phone  string
	Age    int
	Gender string
}

type DoctorCore struct {
	ID        int
	Room      RoomCore
	Specialty string
	Email     string
	Name      string
	Phone     string
	Age       int
	Gender    string
}

type RoomCore struct {
	ID    int
	Floor string
	Code  string
}

type IBusiness interface {
	FindWorkSchedules(q ScheduleQuery) ([]WorkScheduleCore, error)
	FindDoctorWorkSchedules(doctorId int, q ScheduleQuery) ([]WorkScheduleCore, error)
	FindNurseWorkSchedules(nurseId int, q ScheduleQuery) ([]WorkScheduleCore, error)
	CreateWorkSchedule(workSchedule WorkScheduleCore, q ScheduleQuery) error // GENERATE LIST
	EditWorkSchedule(workSchedule WorkScheduleCore) error
	RemoveWorkScheduleById(workScheduleId int) error
	RemoveDoctorFutureWorkSchedules(doctorId int) error
	RemoveNurseFromNextWorkSchedules(nurseId int) error

	FindOutpatietns(q ScheduleQuery) ([]OutpatientCore, error)
	FindOutpatietnsByWorkScheduleId(workScheduleId int) (WorkScheduleCore, error)
	FindOutpatientsByPatientId(patientId int, q ScheduleQuery) ([]OutpatientCore, error)
	FindOutpatientById(outpatientId int) (OutpatientCore, error)
	CreateOutpatient(outpatient OutpatientCore) error

	EditOutpatient(outpatient OutpatientCore) error // ONLY EDIT COMPLAINT
	ExamineOutpatient(outpatientId int, userId int, role string) error
	FinishOutpatient(outpatient OutpatientCore, userId int, role string) error // UpdateOutpatient + InsertPrescriptions
	CancelOutpatient(outpatientId int, userId int, role string) error

	RemoveOutpatientById(outpatientId int) error
	RemovePatientWaitingOutpatients(patientId int) error
}

type IData interface {
	SelectWorkSchedules(q ScheduleQuery) ([]WorkScheduleCore, error)
	SelectCountWorkSchedulesWaitings(ids []int) (map[int]int, error)
	SelectWorkScheduleById(workScheduleId int) (WorkScheduleCore, error)
	SelectWorkSchedulesByDoctorId(doctorId int, q ScheduleQuery) ([]WorkScheduleCore, error)
	SelectWorkSchedulesByNurseId(nurseId int, q ScheduleQuery) ([]WorkScheduleCore, error)
	InsertWorkSchedules(workSchedules []WorkScheduleCore) error
	UpdateWorkSchedule(workSchedule WorkScheduleCore) error
	DeleteWorkScheduleById(workScheduleId int) error // also remove outpatient schedules
	DeleteWorkSchedulesByDoctorId(doctorId int, q ScheduleQuery) error
	DeleteNurseFromWorkSchedules(nurseId int, q ScheduleQuery) error

	SelectOutpatients(q ScheduleQuery) ([]OutpatientCore, error)
	SelectOutpatientsByWorkScheduleId(workScheduleId int) (WorkScheduleCore, error)
	SelectOutpatientsByPatientId(patientId int, q ScheduleQuery) ([]OutpatientCore, error)
	SelectOutpatientById(outpatientId int) (OutpatientCore, error)
	InsertOutpatient(outpatient OutpatientCore) error
	UpdateOutpatient(outpatient OutpatientCore) error
	DeleteWaitingOutpatientsByPatientId(patientId int) error
	DeleteOutpatientById(outpatientId int) error

	// maybe needed in future
	// ? UpdatePresctiption(PrescriptionCore) error
	// ? DeleltePrescriptionById(int) error
}
