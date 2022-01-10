package data

import (
	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"gorm.io/gorm"
)

type mySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepo(db *gorm.DB) schedules.IData {
	return &mySQLRepository{db}
}

func (r *mySQLRepository) SelectWorkSchedules(q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.data.SelectWorkSchedules"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ws := []WorkSchedule{}
	err := r.db.
		Where("date BETWEEN ? AND ?", q.StartDate, q.EndDate).
		Limit(q.Limit).
		Find(&ws).
		Error

	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(err, op, errMsg, errors.KindServerError)
	}
	return toSliceWorkScheduleCore(ws), nil
}

func (r *mySQLRepository) SelectCountWorkSchedulesWaitings(ids []int) (map[int]int, error) {
	const op errors.Op = "schedules.data.SelectCountWorkSchedulesWaitings"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	waitings := []TotalWaiting{}
	query := `
		SELECT w.id, COUNT(w.id) AS total FROM work_schedules w
		JOIN outpatients o 
		ON (
			o.work_schedule_id = w.id AND 
			o.deleted_at IS NULL AND
			w.deleted_at IS NULL AND
			AND o.status = ? AND 
			w.id IN (?)
		)
		GROUP BY w.id
	`

	err := r.db.Raw(query, schedules.StatusWaiting, ids).Scan(&waitings).Error
	if err != nil {
		return map[int]int{}, errors.E(err, op, errMsg, errors.KindServerError)
	}

	result := make(map[int]int)
	for _, w := range waitings {
		result[w.ID] = w.Total
	}
	return result, nil
}

func (r *mySQLRepository) SelectWorkScheduleById(workScheduleId int) (schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.data.SelectWorkScheduleById"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ws := WorkSchedule{}
	err := r.db.First(&ws, workScheduleId).Error
	if err != nil {
		kind := errors.KindServerError
		switch err {
		case gorm.ErrRecordNotFound:
			kind = errors.KindNotFound
			errMsg = "Work schedule not found"
		}
		return schedules.WorkScheduleCore{}, errors.E(err, op, errMsg, kind)
	}

	return ws.toWorkScheduleCore(), nil
}

func (r *mySQLRepository) SelectWorkSchedulesByDoctorId(doctorId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.data.SelectWorkSchedulesByDoctorId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ws := []WorkSchedule{}
	err := r.db.
		Where("doctor_id = ? AND (date BETWEEN ? AND ?) ", doctorId, q.StartDate, q.EndDate).
		Limit(q.Limit).
		Find(&ws).
		Error

	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(err, op, errMsg, errors.KindServerError)
	}

	return toSliceWorkScheduleCore(ws), nil
}

func (r *mySQLRepository) SelectWorkSchedulesByNurseId(nurseId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.data.SelectWorkSchedulesByNurseId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ws := []WorkSchedule{}
	err := r.db.
		Where("nurse_id = ? AND (date BETWEEN ? AND ?) ", nurseId, q.StartDate, q.EndDate).
		Limit(q.Limit).
		Find(&ws).
		Error

	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(err, op, errMsg, errors.KindServerError)
	}

	return toSliceWorkScheduleCore(ws), nil
}

func (r *mySQLRepository) InsertWorkSchedules(workSchedules []schedules.WorkScheduleCore) error {
	const op errors.Op = "schedules.data.InsertWorkSchedules"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ws := make([]WorkSchedule, len(workSchedules))
	for i, w := range workSchedules {
		start, err := NewMyTime(w.StartTime)
		if err != nil {
			errMsg = "Invalid time format"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
		}
		end, err := NewMyTime(w.EndTime)
		if err != nil {
			errMsg = "Invalid time format"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
		}

		ws[i] = WorkSchedule{
			DoctorID:  w.Doctor.ID,
			NurseID:   w.Nurse.ID,
			Group:     w.Group,
			Date:      w.Date,
			StartTime: start,
			EndTime:   end,
		}
	}

	err := r.db.Create(&ws).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) UpdateWorkSchedule(workSchedule schedules.WorkScheduleCore) error {
	const op errors.Op = "schedules.data.UpdateWorkSchedule"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	start, err := NewMyTime(workSchedule.StartTime)
	if err != nil {
		errMsg = "Invalid time format"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}
	end, err := NewMyTime(workSchedule.EndTime)
	if err != nil {
		errMsg = "Invalid time format"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	updatedWorkSchedule := WorkSchedule{
		Model: gorm.Model{
			ID:        uint(workSchedule.ID),
			CreatedAt: workSchedule.CreatedAt,
		},
		DoctorID:  workSchedule.Doctor.ID,
		NurseID:   workSchedule.Nurse.ID,
		Group:     workSchedule.Group,
		Date:      workSchedule.Date,
		StartTime: start,
		EndTime:   end,
	}

	err = r.db.Save(&updatedWorkSchedule).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) DeleteWorkScheduleById(workScheduleId int) error {
	// also remove outpatient schedules
	const op errors.Op = "schedules.data.DeleteWorkScheduleById"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	trasaction := func(tx *gorm.DB) error {
		os := []Outpatient{}

		err := tx.Where("work_schedule_id = ?", workScheduleId).Find(&os).Error
		if err != nil {
			return err
		}

		outpatientIds := make([]uint, len(os))
		for i := range os {
			outpatientIds[i] = os[i].ID
		}

		err = tx.Where("outpatient_id IN (?)", outpatientIds).Delete(&Prescription{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("work_schedule_id = ?", workScheduleId).Delete(&Outpatient{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&WorkSchedule{}, workScheduleId).Error
		if err != nil {
			return err
		}

		return nil
	}

	err := r.db.Transaction(trasaction)
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) DeleteWorkSchedulesByDoctorId(doctorId int, q schedules.ScheduleQuery) error {
	const op errors.Op = "schedules.data.DeleteWorkSchedulesByDoctorId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	err := r.db.
		Where("doctor_id = ? AND (date BETWEEN ? AND ?)", doctorId, q.StartDate, q.EndDate).
		Delete(&WorkSchedule{}).
		Error

	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) DeleteNurseFromWorkSchedules(nurseId int, q schedules.ScheduleQuery) error {
	const op errors.Op = "schedules.data.DeleteWorkSchedulesByNurseId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	query := `
		UPDATE work_schedules
		SET nurse_id = NULL
		WHERE nurse_id = ? AND (date BETWEEN ? AND ?)
	`
	err := r.db.Exec(query, nurseId, q.StartDate, q.EndDate).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) SelectOutpatients(q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.data.SelectOutpatients"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	query := `
	SELECT 
		outpatients.id, outpatients.created_at, outpatients.updated_at, outpatients.deleted_at, outpatients.work_schedule_id, outpatients.patient_id, 
		outpatients.complaint, outpatients.status, outpatients.start_time, outpatients.end_time, 
		WorkSchedule.id AS WorkSchedule__id, WorkSchedule.created_at AS WorkSchedule__created_at, WorkSchedule.updated_at AS WorkSchedule__updated_at, 
		WorkSchedule.deleted_at AS WorkSchedule__deleted_at, WorkSchedule.doctor_id AS WorkSchedule__doctor_id, 
		WorkSchedule.nurse_id AS WorkSchedule__nurse_id, WorkSchedule.group AS WorkSchedule__group, WorkSchedule.date AS WorkSchedule__date, 
		WorkSchedule.start_time AS WorkSchedule__start_time, WorkSchedule.end_time AS WorkSchedule__end_time FROM outpatients 
	JOIN work_schedules WorkSchedule 
	ON (
		outpatients.work_schedule_id = WorkSchedule.id AND 
		outpatients.deleted_at IS NULL AND 
		WorkSchedule.deleted_at IS NULL AND 
		(WorkSchedule.date BETWEEN ? AND ?)
	)
	LIMIT ?
	`
	os := []Outpatient{}
	err := r.db.Raw(query, q.StartDate, q.EndDate, q.Limit).Scan(&os).Error
	/* OLD WAY
	err := r.db.
		Joins("WorkSchedule").
		Where("date BETWEEN ? AND ?", q.StartDate, q.EndDate).
		Limit(q.Limit).
		Find(&os).
		Error
	*/

	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op, errMsg, errors.KindServerError)
	}

	return toSliceOutpatientCore(os), nil
}

func (r *mySQLRepository) SelectOutpatientsByPatientId(patientId int, q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.data.SelectOutpatientsByPatientId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	query := `
	SELECT 
		outpatients.id, outpatients.created_at, outpatients.updated_at, outpatients.deleted_at, outpatients.work_schedule_id, outpatients.patient_id, 
		outpatients.complaint, outpatients.status, outpatients.start_time, outpatients.end_time, 
		WorkSchedule.id AS WorkSchedule__id, WorkSchedule.created_at AS WorkSchedule__created_at, WorkSchedule.updated_at AS WorkSchedule__updated_at, 
		WorkSchedule.deleted_at AS WorkSchedule__deleted_at, WorkSchedule.doctor_id AS WorkSchedule__doctor_id, 
		WorkSchedule.nurse_id AS WorkSchedule__nurse_id, WorkSchedule.group AS WorkSchedule__group, WorkSchedule.date AS WorkSchedule__date, 
		WorkSchedule.start_time AS WorkSchedule__start_time, WorkSchedule.end_time AS WorkSchedule__end_time FROM outpatients 
	JOIN work_schedules WorkSchedule 
	ON (
		outpatients.work_schedule_id = WorkSchedule.id AND 
		outpatients.deleted_at IS NULL AND 
		WorkSchedule.deleted_at IS NULL AND 
		patient_id = ? AND 
		(WorkSchedule.date BETWEEN ? AND ?)
	)
	LIMIT ?
	`
	os := []Outpatient{}
	err := r.db.Raw(query, patientId, q.StartDate, q.EndDate, q.Limit).Scan(&os).Error

	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op, errMsg, errors.KindServerError)
	}

	return toSliceOutpatientCore(os), nil
}

func (r *mySQLRepository) SelectOutpatientsByWorkScheduleId(workScheduleId int) (schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.data.SelectOutpatientsByWorkScheduleId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	w := WorkSchedule{}

	// err := r.db.Where("work_schedule_id = ?", workScheduleId).Find(&os).Error
	err := r.db.Preload("Outpatients").First(&w, workScheduleId).Error
	if err != nil {
		kind := errors.KindServerError
		switch err {
		case gorm.ErrRecordNotFound:
			errMsg = "Work schedule not found"
			kind = errors.KindNotFound
		}
		return schedules.WorkScheduleCore{}, errors.E(err, op, errMsg, kind)
	}

	return w.toWorkScheduleCore(), nil
}

func (r *mySQLRepository) SelectOutpatientById(outpatientId int) (schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.data.SelectOutpatientById"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	o := Outpatient{}
	err := r.db.
		Preload("WorkSchedule").
		Preload("Prescriptions").
		First(&o, outpatientId).
		Error

	if err != nil {
		kind := errors.KindServerError
		switch err {
		case gorm.ErrRecordNotFound:
			kind = errors.KindNotFound
			errMsg = "Outpatient not found"
		}
		return schedules.OutpatientCore{}, errors.E(err, op, errMsg, kind)
	}

	return o.toOutpatientCore(), nil
}

func (r *mySQLRepository) InsertOutpatient(outpatient schedules.OutpatientCore) error {
	const op errors.Op = "schedules.data.InsertOutpatient"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	newOutpatient := Outpatient{
		WorkScheduleID: uint(outpatient.WorkSchedule.ID),
		PatientID:      outpatient.Patient.ID,
		Complaint:      outpatient.Complaint,
		Status:         outpatient.Status,
	}

	err := r.db.Create(&newOutpatient).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) UpdateOutpatient(outpatient schedules.OutpatientCore) error {
	const op errors.Op = "schedules.data.UpdateOutpatient"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	ps := make([]Prescription, len(outpatient.Prescriptions))
	for i := range outpatient.Prescriptions {
		ps[i] = Prescription{
			OutpatientID: uint(outpatient.ID),
			Medicine:     outpatient.Prescriptions[i].Medicine,
			Instruction:  outpatient.Prescriptions[i].Instruction,
		}
	}

	start, err := NewMyTime(outpatient.StartTime)
	if err != nil {
		errMsg = "Invalid time format"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}
	end, err := NewMyTime(outpatient.EndTime)
	if err != nil {
		errMsg = "Invalid time format"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindBadRequest)
	}

	updatedOutpatient := Outpatient{
		Model:          gorm.Model{ID: uint(outpatient.ID), CreatedAt: outpatient.CreatedAt},
		WorkScheduleID: uint(outpatient.WorkSchedule.ID),
		PatientID:      outpatient.Patient.ID,
		Complaint:      outpatient.Complaint,
		Status:         outpatient.Status,
		StartTime:      start,
		EndTime:        end,
		Prescriptions:  ps,
	}

	err = r.db.Save(&updatedOutpatient).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) DeleteWaitingOutpatientsByPatientId(patientId int) error {
	const op errors.Op = "schedules.data.DeleteWaitingOutpatientsByPatientId"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	err := r.db.Where("patient_id = ? AND status = ?", patientId, schedules.StatusWaiting).Delete(Outpatient{}).Error
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}

func (r *mySQLRepository) DeleteOutpatientById(outpatientId int) error {
	const op errors.Op = "schedules.data.DeleteOutpatientById"
	var errMsg errors.ErrClientMessage = "Something went wrong"

	outpatientDeletion := func(tx *gorm.DB) error {
		err := tx.Where("outpatient_id = ?", outpatientId).Delete(&Prescription{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&Outpatient{}, outpatientId).Error
		if err != nil {
			return err
		}

		return nil
	}

	err := r.db.Transaction(outpatientDeletion)
	if err != nil {
		return errors.E(err, op, errMsg, errors.KindServerError)
	}

	return nil
}
