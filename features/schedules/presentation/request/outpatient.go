package request

import "github.com/final-project-alterra/hospital-management-system-api/features/schedules"

type CreateOutpatientRequest struct {
	WorkScheduleID int    `json:"workScheduleId" validate:"gt=0"`
	PatientID      int    `json:"patientId" validate:"gt=0"`
	Complaint      string `json:"complaint" validate:"required"`
}

func (o CreateOutpatientRequest) ToOutpatientCore() schedules.OutpatientCore {
	core := schedules.OutpatientCore{}
	core.WorkSchedule.ID = o.WorkScheduleID
	core.Patient.ID = o.PatientID
	core.Complaint = o.Complaint

	return core
}

type UpdateOutpatientRequest struct {
	ID        int    `json:"id" validate:"gt=0"`
	Complaint string `json:"complaint" validate:"required"`
}

func (o UpdateOutpatientRequest) ToOutpatientCore() schedules.OutpatientCore {
	core := schedules.OutpatientCore{}
	core.ID = o.ID
	core.Complaint = o.Complaint

	return core
}

type FinishOutpatientRequest struct {
	ID            int                   `json:"id" validate:"gt=0"`
	Diagnosis     string                `json:"diagnosis" validate:"required"`
	Prescriptions []PrescriptionRequest `json:"prescriptions" validate:"required"`
}

func (o FinishOutpatientRequest) ToOutpatientCore() schedules.OutpatientCore {
	pc := make([]schedules.PrescriptionCore, len(o.Prescriptions))

	for i, p := range o.Prescriptions {
		pc[i] = p.ToPrescriptionCore()
	}

	return schedules.OutpatientCore{
		ID:            o.ID,
		Diagnosis:     o.Diagnosis,
		Prescriptions: pc,
	}
}

type PrescriptionRequest struct {
	Medicine    string `json:"medicine" validate:"required"`
	Instruction string `json:"instruction" validate:"required"`
}

func (p PrescriptionRequest) ToPrescriptionCore() schedules.PrescriptionCore {
	return schedules.PrescriptionCore{
		Medicine:    p.Medicine,
		Instruction: p.Instruction,
	}
}

type ExamineOutpatientRequest struct {
	ID int `json:"id" validate:"gt=0"`
}

type CancelOutpatientRequest struct {
	ID int `json:"id" validate:"gt=0"`
}
