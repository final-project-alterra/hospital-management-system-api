package response

import "github.com/final-project-alterra/hospital-management-system-api/features/schedules"

type OutpatientResponse struct {
	ID        int                `json:"id"`
	Status    int                `json:"status"`
	Date      string             `json:"date"`
	StartTime string             `json:"startTime"`
	EndTime   string             `json:"endTime"`
	Complaint string             `json:"complaint"`
	Patient   Outpatient_Patient `json:"patient"`
	Doctor    Outpatient_Doctor  `json:"doctor"`
	Nurse     Outpatient_Nurse   `json:"nurse"`
}

type PatientOutpatientResponse struct {
	ID        int    `json:"id"`
	Status    int    `json:"status"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Complaint string `json:"complaint"`

	Doctor Outpatient_Doctor `json:"doctor"`
	Nurse  Outpatient_Nurse  `json:"nurse"`
}

type WorkScheduleOutPatientResponse struct {
	ID        int               `json:"id"`
	Date      string            `json:"date"`
	StartTime string            `json:"startTime"`
	EndTime   string            `json:"endTime"`
	Doctor    Outpatient_Doctor `json:"doctor"`
	Nurse     Outpatient_Nurse  `json:"nurse"`

	Outpatients []Outpatient_WorkScheduleOutPatient_Outpatient `json:"outpatients"`
}

type OutpatientDetailResponse struct {
	ID           int                    `json:"id"`
	Status       int                    `json:"status"`
	Date         string                 `json:"date"`
	StartTime    string                 `json:"startTime"`
	EndTime      string                 `json:"endTime"`
	Complaint    string                 `json:"complaint"`
	Patient      Outpatient_Patient     `json:"patient"`
	Doctor       Outpatient_Doctor      `json:"doctor"`
	Nurse        Outpatient_Nurse       `json:"nurse"`
	Prescription []PrescriptionResponse `json:"prescription"`
}

type PrescriptionResponse struct {
	ID          int    `json:"id"`
	Medicine    string `json:"medicine"`
	Instruction string `json:"instruction"`
}

func Outpatient(o schedules.OutpatientCore) OutpatientResponse {
	return OutpatientResponse{
		ID:        o.ID,
		Status:    o.Status,
		Date:      o.WorkSchedule.Date,
		StartTime: o.StartTime,
		EndTime:   o.EndTime,
		Complaint: o.Complaint,

		Patient: Outpatient_Patient{}.FromCore(o.Patient),
		Doctor:  Outpatient_Doctor{}.FromCore(o.WorkSchedule.Doctor),
		Nurse:   Outpatient_Nurse{}.FromCore(o.WorkSchedule.Nurse),
	}
}

func PatientOutpatient(o schedules.OutpatientCore) PatientOutpatientResponse {
	return PatientOutpatientResponse{
		ID:        o.ID,
		Status:    o.Status,
		Date:      o.WorkSchedule.Date,
		StartTime: o.StartTime,
		EndTime:   o.EndTime,
		Complaint: o.Complaint,

		Doctor: Outpatient_Doctor{}.FromCore(o.WorkSchedule.Doctor),
		Nurse:  Outpatient_Nurse{}.FromCore(o.WorkSchedule.Nurse),
	}
}

func WorkScheduleOutpatient(w schedules.WorkScheduleCore) WorkScheduleOutPatientResponse {
	o := make([]Outpatient_WorkScheduleOutPatient_Outpatient, len(w.Outpatients))

	for i := range w.Outpatients {
		o[i] = Outpatient_WorkScheduleOutPatient_Outpatient{}.FromCore(w.Outpatients[i])
	}

	return WorkScheduleOutPatientResponse{
		ID:          w.ID,
		StartTime:   w.StartTime,
		EndTime:     w.EndTime,
		Date:        w.Date,
		Doctor:      Outpatient_Doctor{}.FromCore(w.Doctor),
		Nurse:       Outpatient_Nurse{}.FromCore(w.Nurse),
		Outpatients: o,
	}
}

func OutpatientDetail(o schedules.OutpatientCore) OutpatientDetailResponse {
	return OutpatientDetailResponse{
		ID:        o.ID,
		Status:    o.Status,
		Date:      o.WorkSchedule.Date,
		StartTime: o.StartTime,
		EndTime:   o.EndTime,
		Complaint: o.Complaint,

		Patient: Outpatient_Patient{}.FromCore(o.Patient),
		Doctor:  Outpatient_Doctor{}.FromCore(o.WorkSchedule.Doctor),
		Nurse:   Outpatient_Nurse{}.FromCore(o.WorkSchedule.Nurse),

		Prescription: ListPrescription(o.Prescriptions),
	}
}

func Prescription(p schedules.PrescriptionCore) PrescriptionResponse {
	return PrescriptionResponse{
		ID:          p.ID,
		Medicine:    p.Medicine,
		Instruction: p.Instruction,
	}
}

/* List */
func ListOutpatients(o []schedules.OutpatientCore) []OutpatientResponse {
	resp := make([]OutpatientResponse, len(o))

	for i := range o {
		resp[i] = Outpatient(o[i])
	}

	return resp
}

func ListPatientOutpatients(o []schedules.OutpatientCore) []PatientOutpatientResponse {
	resp := make([]PatientOutpatientResponse, len(o))

	for i := range o {
		resp[i] = PatientOutpatient(o[i])
	}

	return resp
}

func ListPrescription(ps []schedules.PrescriptionCore) []PrescriptionResponse {
	resp := make([]PrescriptionResponse, len(ps))

	for i := range ps {
		resp[i] = Prescription(ps[i])
	}

	return resp
}

/* Nested struct for outpatients */
type Outpatient_Patient struct {
	ID    int    `json:"id"`
	NIK   string `json:"nik"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Outpatient_Doctor struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Specialty string `json:"specialty"`
}

type Outpatient_Nurse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Outpatient_WorkScheduleOutPatient_Outpatient struct {
	ID        int                `json:"id"`
	Status    int                `json:"status"`
	StartTime string             `json:"startTime"`
	EndTime   string             `json:"endTime"`
	Complaint string             `json:"complaint"`
	Patient   Outpatient_Patient `json:"patient"`
}

func (p Outpatient_Patient) FromCore(c schedules.PatientCore) Outpatient_Patient {
	p.ID = c.ID
	p.NIK = c.NIK
	p.Name = c.Name
	p.Phone = c.Phone

	return p
}

func (d Outpatient_Doctor) FromCore(c schedules.DoctorCore) Outpatient_Doctor {
	d.ID = c.ID
	d.Email = c.Email
	d.Name = c.Name
	d.Specialty = c.Specialty

	return d
}

func (n Outpatient_Nurse) FromCore(c schedules.NurseCore) Outpatient_Nurse {
	n.ID = c.ID
	n.Email = c.Email
	n.Name = c.Name

	return n
}

func (wo Outpatient_WorkScheduleOutPatient_Outpatient) FromCore(o schedules.OutpatientCore) Outpatient_WorkScheduleOutPatient_Outpatient {
	return Outpatient_WorkScheduleOutPatient_Outpatient{
		ID:        o.ID,
		Status:    o.Status,
		StartTime: o.StartTime,
		EndTime:   o.EndTime,
		Complaint: o.Complaint,

		Patient: Outpatient_Patient{}.FromCore(o.Patient),
	}
}
