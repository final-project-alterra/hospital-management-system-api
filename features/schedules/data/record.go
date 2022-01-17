package data

import (
	"strings"

	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"gorm.io/gorm"
)

type WorkSchedule struct {
	gorm.Model
	DoctorID    int
	NurseID     int
	Group       string `gorm:"type:varchar(64)"`
	Date        string `gorm:"type:date"`
	StartTime   MyTime `gorm:"default:null"`
	EndTime     MyTime `gorm:"default:null"`
	Outpatients []Outpatient
}

type Outpatient struct {
	gorm.Model
	WorkScheduleID uint
	WorkSchedule   WorkSchedule

	PatientID     int
	Complaint     string
	Status        int
	StartTime     MyTime `gorm:"default:null"`
	EndTime       MyTime `gorm:"default:null"`
	Prescriptions []Prescription
}

type Prescription struct {
	gorm.Model
	OutpatientID uint

	Medicine    string `gorm:"type:varchar(64)"`
	Instruction string
}

// SELECT id, COUNT(*) FROM work_schedules GROUP BY id HAVING COUNT(*) > 1;
type TotalWaiting struct {
	ID    int // ID of the work schedule
	Total int
}

func (w *WorkSchedule) toWorkScheduleCore() schedules.WorkScheduleCore {

	return schedules.WorkScheduleCore{
		ID:          int(w.ID),
		Group:       w.Group,
		Date:        strings.Split(w.Date, "T")[0],
		StartTime:   w.StartTime.String(),
		EndTime:     w.EndTime.String(),
		Doctor:      schedules.DoctorCore{ID: w.DoctorID},
		Nurse:       schedules.NurseCore{ID: w.NurseID},
		Outpatients: toSliceOutpatientCore(w.Outpatients),
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

func (o *Outpatient) toOutpatientCore() schedules.OutpatientCore {
	return schedules.OutpatientCore{
		ID:            int(o.ID),
		Complaint:     o.Complaint,
		Status:        o.Status,
		StartTime:     o.StartTime.String(),
		EndTime:       o.EndTime.String(),
		Patient:       schedules.PatientCore{ID: o.PatientID},
		WorkSchedule:  o.WorkSchedule.toWorkScheduleCore(),
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
		Prescriptions: toSlicePrescriptionCore(o.Prescriptions),
	}
}

func (p *Prescription) toPrescriptionCore() schedules.PrescriptionCore {
	return schedules.PrescriptionCore{
		ID:          int(p.ID),
		Medicine:    p.Medicine,
		Instruction: p.Instruction,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toSliceWorkScheduleCore(ws []WorkSchedule) []schedules.WorkScheduleCore {
	wsc := make([]schedules.WorkScheduleCore, len(ws))
	for i := range ws {
		wsc[i] = ws[i].toWorkScheduleCore()
	}
	return wsc
}

func toSliceOutpatientCore(o []Outpatient) []schedules.OutpatientCore {
	oc := make([]schedules.OutpatientCore, len(o))
	for i := range o {
		oc[i] = o[i].toOutpatientCore()
	}
	return oc
}

func toSlicePrescriptionCore(p []Prescription) []schedules.PrescriptionCore {
	pc := make([]schedules.PrescriptionCore, len(p))
	for i := range p {
		pc[i] = p[i].toPrescriptionCore()
	}
	return pc
}
