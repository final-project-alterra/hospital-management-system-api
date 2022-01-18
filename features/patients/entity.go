package patients

import "time"

type PatientCore struct {
	ID        int
	CreatedBy int
	UpdatedBy int
	NIK       string
	Name      string
	BirthDate string
	Phone     string
	Address   string
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBusiness interface {
	FindPatients() ([]PatientCore, error)
	FindPatientsByIds(ids []int) ([]PatientCore, error)
	FindPatientById(id int) (PatientCore, error)
	CreatePatient(patient PatientCore) error
	EditPatient(patient PatientCore) error
	RemovePatientById(id int, updatedBy int) error
}

type IData interface {
	SelectPatients() ([]PatientCore, error)
	SelectPatientsByIds(ids []int) ([]PatientCore, error)
	SelectPatientById(id int) (PatientCore, error)
	SelectPatientByNIK(nik string) (PatientCore, error)
	InsertPatient(patient PatientCore) error
	UpdatePatient(patient PatientCore) error
	DeletePatientById(id int, updatedBy int) error
}
