package business

import (
	"time"

	"github.com/final-project-alterra/hospital-management-system-api/errors"
	"github.com/final-project-alterra/hospital-management-system-api/features/doctors"
	"github.com/final-project-alterra/hospital-management-system-api/features/nurses"
	"github.com/final-project-alterra/hospital-management-system-api/features/patients"
	"github.com/final-project-alterra/hospital-management-system-api/features/schedules"
	"github.com/google/uuid"
)

type scheduleBusiness struct {
	data            schedules.IData
	doctorBusiness  doctors.IBusiness
	nurseBusiness   nurses.IBusiness
	patientBusiness patients.IBusiness
}

func (s *scheduleBusiness) FindWorkSchedules(q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.business.FindWorkSchedules"

	schedulesData, err := s.data.SelectWorkSchedules(q)
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	scheduleIds := make([]int, len(schedulesData))
	for i := range schedulesData {
		scheduleIds[i] = schedulesData[i].ID
	}

	waitingMap, err := s.data.SelectCountWorkSchedulesWaitings(scheduleIds)
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	doctorsMap, err := s.findDoctorsData(s.getUniqueDoctorIds(schedulesData))
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	nursesMap, err := s.findNursesData(s.getUniqueNurseIds(schedulesData))
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	for i := range schedulesData {
		ID := schedulesData[i].ID
		doctorID := schedulesData[i].Doctor.ID
		nurseID := schedulesData[i].Nurse.ID

		schedulesData[i].TotalWaiting = waitingMap[ID]
		schedulesData[i].Doctor = doctorsMap[doctorID]
		schedulesData[i].Nurse = nursesMap[nurseID]
	}

	return schedulesData, nil
}

func (s *scheduleBusiness) FindDoctorWorkSchedules(doctorId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.business.FindDoctorWorkSchedules"

	schedulesData, err := s.data.SelectWorkSchedulesByDoctorId(doctorId, q)
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	nursesMap, err := s.findNursesData(s.getUniqueNurseIds(schedulesData))
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	for i := range schedulesData {
		nurseID := schedulesData[i].Nurse.ID
		schedulesData[i].Nurse = nursesMap[nurseID]
	}

	return schedulesData, nil
}

func (s *scheduleBusiness) FindNurseWorkSchedules(nurseId int, q schedules.ScheduleQuery) ([]schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.business.FindNurseWorkSchedules"

	schedulesData, err := s.data.SelectWorkSchedulesByNurseId(nurseId, q)
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	doctorsMap, err := s.findDoctorsData(s.getUniqueDoctorIds(schedulesData))
	if err != nil {
		return []schedules.WorkScheduleCore{}, errors.E(err, op)
	}

	for i := range schedulesData {
		doctorID := schedulesData[i].Doctor.ID
		schedulesData[i].Doctor = doctorsMap[doctorID]
	}

	return schedulesData, nil
}

func (s *scheduleBusiness) CreateWorkSchedule(workSchedule schedules.WorkScheduleCore, q schedules.ScheduleQuery) error { // GENERATE LIST
	const op errors.Op = "schedules.business.CreateWorkSchedule"
	var errMesage errors.ErrClientMessage

	var dates []string

	_, err := s.doctorBusiness.FindDoctorById(workSchedule.Doctor.ID)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = s.nurseBusiness.FindNurseById(workSchedule.Nurse.ID)
	if err != nil {
		return errors.E(err, op)
	}

	switch q.Repeat {
	case schedules.RepeatNoRepeat:
		dates = []string{q.StartDate}
	case schedules.RepeatDaily:
		dates, err = s.repeatEveryDay(q.StartDate, q.EndDate)
	case schedules.RepeatWeekly:
		dates, err = s.repeatEveryWeek(q.StartDate, q.EndDate)
	case schedules.RepeatMonthly:
		dates, err = s.repeatEveryMonthSameDay(q.StartDate, q.EndDate)
	default:
		errMesage = "Invalid repeat type"
		return errors.E(errors.New("Invalid repeat type"), op, errMesage, errors.KindBadRequest)
	}

	if err != nil {
		return errors.E(err, op)
	}

	// ! Potential panic
	group := uuid.New().String()

	newSchedules := make([]schedules.WorkScheduleCore, len(dates))
	for i := range dates {
		newSchedule := workSchedule
		newSchedule.Date = dates[i]
		newSchedule.Group = group

		newSchedules[i] = newSchedule
	}

	err = s.data.InsertWorkSchedules(newSchedules)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (s *scheduleBusiness) EditWorkSchedule(workSchedule schedules.WorkScheduleCore) error {
	const op errors.Op = "schedules.business.EditWorkSchedule"

	existingSchedules, err := s.data.SelectWorkScheduleById(workSchedule.ID)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = s.doctorBusiness.FindDoctorById(workSchedule.Doctor.ID)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = s.nurseBusiness.FindNurseById(workSchedule.Nurse.ID)
	if err != nil {
		return errors.E(err, op)
	}

	existingSchedules.Doctor.ID = workSchedule.Doctor.ID
	existingSchedules.Nurse.ID = workSchedule.Nurse.ID
	existingSchedules.Date = workSchedule.Date
	existingSchedules.StartTime = workSchedule.StartTime
	existingSchedules.EndTime = workSchedule.EndTime

	err = s.data.UpdateWorkSchedule(existingSchedules)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (s *scheduleBusiness) RemoveWorkScheduleById(workScheduleId int) error {
	const op errors.Op = "schedules.business.RemoveWorkScheduleById"

	err := s.data.DeleteWorkScheduleById(workScheduleId)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (s *scheduleBusiness) RemoveDoctorFutureWorkSchedules(doctorId int) error {
	const op errors.Op = "schedules.business.RemoveDoctorFutureWorkSchedules"

	q := schedules.ScheduleQuery{
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().AddDate(100, 0, 0).Format("2006-01-02"),
	}
	err := s.data.DeleteWorkSchedulesByDoctorId(doctorId, q)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) RemoveNurseFromNextWorkSchedules(nurseId int) error {
	const op errors.Op = "schedules.business.RemoveNurseFromNextWorkSchedules"

	q := schedules.ScheduleQuery{
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().AddDate(100, 0, 0).Format("2006-01-02"),
	}
	err := s.data.DeleteNurseFromWorkSchedules(nurseId, q)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) FindOutpatietns(q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.business.FindOutpatietns"

	outpatientsData, err := s.data.SelectOutpatients(q)
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(op, err)
	}

	patientsMap, err := s.findPatientData(s.getUniquePatientIds(outpatientsData))
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op)
	}

	uniqueSchedules := s.getUniqueSchedules(outpatientsData)

	doctorsMap, err := s.findDoctorsData(s.getUniqueDoctorIds(uniqueSchedules))
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op)
	}

	nurseMap, err := s.findNursesData(s.getUniqueNurseIds(uniqueSchedules))
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op)
	}

	for i := range outpatientsData {
		patientID := outpatientsData[i].Patient.ID
		doctorID := outpatientsData[i].WorkSchedule.Doctor.ID
		nurseID := outpatientsData[i].WorkSchedule.Nurse.ID

		outpatientsData[i].Patient = patientsMap[patientID]
		outpatientsData[i].WorkSchedule.Doctor = doctorsMap[doctorID]
		outpatientsData[i].WorkSchedule.Nurse = nurseMap[nurseID]
	}

	return outpatientsData, nil
}

func (s *scheduleBusiness) FindOutpatietnsByWorkScheduleId(workScheduleId int) (schedules.WorkScheduleCore, error) {
	const op errors.Op = "schedules.business.FindOutpatietnsByWorkScheduleId"

	workSchedule, err := s.data.SelectOutpatientsByWorkScheduleId(workScheduleId)
	if err != nil {
		return schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	doctor, err := s.doctorBusiness.FindDoctorById(workSchedule.Doctor.ID)
	if err != nil {
		return schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	nurse, err := s.nurseBusiness.FindNurseById(workSchedule.Nurse.ID)
	if err != nil {
		return schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	patientsMap, err := s.findPatientData(s.getUniquePatientIds(workSchedule.Outpatients))
	if err != nil {
		return schedules.WorkScheduleCore{}, errors.E(op, err)
	}

	doctorCore := schedules.DoctorCore{
		ID:        doctor.ID,
		Name:      doctor.Name,
		Email:     doctor.Email,
		Phone:     doctor.Phone,
		Specialty: doctor.Speciality.Name,
		Age:       doctor.Age,
		Gender:    doctor.Gender,
	}
	nurseCore := schedules.NurseCore{
		ID:     nurse.ID,
		Name:   nurse.Name,
		Email:  nurse.Email,
		Phone:  nurse.Phone,
		Age:    nurse.Age,
		Gender: nurse.Gender,
	}

	workSchedule.Doctor = doctorCore
	workSchedule.Nurse = nurseCore

	for i := range workSchedule.Outpatients {
		patientID := workSchedule.Outpatients[i].Patient.ID
		workSchedule.Outpatients[i].Patient = patientsMap[patientID]
	}

	return workSchedule, nil
}

func (s *scheduleBusiness) FindOutpatientsByPatientId(patientId int, q schedules.ScheduleQuery) ([]schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.business.FindOutpatientsByPatientId"

	outpatientsData, err := s.data.SelectOutpatientsByPatientId(patientId, q)
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(op, err)
	}

	uniqueSchedules := s.getUniqueSchedules(outpatientsData)

	doctorsMap, err := s.findDoctorsData(s.getUniqueDoctorIds(uniqueSchedules))
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op)
	}

	nurseMap, err := s.findNursesData(s.getUniqueNurseIds(uniqueSchedules))
	if err != nil {
		return []schedules.OutpatientCore{}, errors.E(err, op)
	}

	for i := range outpatientsData {
		doctorID := outpatientsData[i].WorkSchedule.Doctor.ID
		nurseID := outpatientsData[i].WorkSchedule.Nurse.ID

		outpatientsData[i].WorkSchedule.Doctor = doctorsMap[doctorID]
		outpatientsData[i].WorkSchedule.Nurse = nurseMap[nurseID]
	}

	return outpatientsData, nil
}

func (s *scheduleBusiness) FindOutpatientById(outpatientId int) (schedules.OutpatientCore, error) {
	const op errors.Op = "schedules.business.FindOutpatientById"

	outpatientData, err := s.data.SelectOutpatientById(outpatientId)
	if err != nil {
		return schedules.OutpatientCore{}, errors.E(err, op)
	}

	patientData, err := s.patientBusiness.FindPatientById(outpatientData.Patient.ID)
	if err != nil {
		return schedules.OutpatientCore{}, errors.E(err, op)
	}

	doctorData, err := s.doctorBusiness.FindDoctorById(outpatientData.WorkSchedule.Doctor.ID)
	if err != nil {
		return schedules.OutpatientCore{}, errors.E(err, op)
	}

	nurseData, err := s.nurseBusiness.FindNurseById(outpatientData.WorkSchedule.Nurse.ID)
	if err != nil {
		return schedules.OutpatientCore{}, errors.E(err, op)
	}

	patient := schedules.PatientCore{
		ID:     patientData.ID,
		NIK:    patientData.NIK,
		Name:   patientData.Name,
		Phone:  patientData.Phone,
		Age:    patientData.Age,
		Gender: patientData.Gender,
	}
	doctor := schedules.DoctorCore{
		ID:        doctorData.ID,
		Email:     doctorData.Email,
		Name:      doctorData.Name,
		Specialty: doctorData.Speciality.Name,
		Phone:     doctorData.Phone,
		Age:       doctorData.Age,
		Gender:    doctorData.Gender,
		Room:      schedules.RoomCore{ID: doctorData.Room.ID, Code: doctorData.Room.Code, Floor: doctorData.Room.Floor},
	}
	nurse := schedules.NurseCore{
		ID:     nurseData.ID,
		Email:  nurseData.Email,
		Name:   nurseData.Name,
		Phone:  nurseData.Phone,
		Age:    nurseData.Age,
		Gender: nurseData.Gender,
	}

	outpatientData.Patient = patient
	outpatientData.WorkSchedule.Doctor = doctor
	outpatientData.WorkSchedule.Nurse = nurse

	return outpatientData, nil
}

func (s *scheduleBusiness) CreateOutpatient(outpatient schedules.OutpatientCore) error {
	const op errors.Op = "schedules.business.CreateOutpatient"

	_, err := s.data.SelectWorkScheduleById(outpatient.WorkSchedule.ID)
	if err != nil {
		return errors.E(err, op)
	}

	_, err = s.patientBusiness.FindPatientById(outpatient.Patient.ID)
	if err != nil {
		return errors.E(err, op)
	}

	outpatient.Status = schedules.StatusWaiting
	err = s.data.InsertOutpatient(outpatient)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (s *scheduleBusiness) EditOutpatient(outpatient schedules.OutpatientCore) error {
	// ONLY EDIT COMPLAINT
	const op errors.Op = "schedules.business.EditOutpatient"

	existingOutpatient, err := s.data.SelectOutpatientById(outpatient.ID)
	if err != nil {
		return errors.E(err, op)
	}

	existingOutpatient.Complaint = outpatient.Complaint
	err = s.data.UpdateOutpatient(existingOutpatient)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) ExamineOutpatient(outpatientId int, userId int, role string) error {
	const op errors.Op = "schedules.business.ExamineOutpatient"
	var errMsg errors.ErrClientMessage

	existingOutpatient, err := s.data.SelectOutpatientById(outpatientId)
	if err != nil {
		return errors.E(err, op)
	}

	if existingOutpatient.Status != schedules.StatusWaiting {
		errMsg = "Cannot examine patient that is not waiting"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}

	switch role {
	case "doctor":
		if userId != existingOutpatient.WorkSchedule.Doctor.ID {
			errMsg = "Only doctor of this outpatient work schedule can examine this outpatient"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
		}

	case "nurse":
		if userId != existingOutpatient.WorkSchedule.Nurse.ID {
			errMsg = "Only nurse of this outpatient work schedule can examine this outpatient"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
		}

	default:
		errMsg = "Only doctor or nurse of this outpatient work schedule can examine this outpatient"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
	}

	workSchedule, err := s.data.SelectOutpatientsByWorkScheduleId(existingOutpatient.WorkSchedule.ID)
	if err != nil {
		return errors.E(err, op)
	}

	for i := range workSchedule.Outpatients {
		if workSchedule.Outpatients[i].Status == schedules.StatusOnprogress {
			errMsg = "There is an ongoing examination"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
		}
	}

	existingOutpatient.Status = schedules.StatusOnprogress
	existingOutpatient.StartTime = time.Now().Format("15:04:05")

	err = s.data.UpdateOutpatient(existingOutpatient)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) FinishOutpatient(outpatient schedules.OutpatientCore, userId int, role string) error {
	const op errors.Op = "schedules.business.FinishOutpatient"
	var errMsg errors.ErrClientMessage

	existingOutpatient, err := s.data.SelectOutpatientById(outpatient.ID)
	if err != nil {
		return errors.E(err, op)
	}

	if existingOutpatient.Status != schedules.StatusOnprogress {
		errMsg = "Outpatient is not on progress"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}

	switch role {
	case "doctor":
		if userId != existingOutpatient.WorkSchedule.Doctor.ID {
			errMsg = "Only doctor of this outpatient work schedule can finish this outpatient"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
		}

	default:
		errMsg = "Only doctor of this outpatient work schedule can finish this outpatient"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
	}

	existingOutpatient.EndTime = time.Now().Format("15:04:05")
	existingOutpatient.Status = schedules.StatusFinished
	existingOutpatient.Prescriptions = outpatient.Prescriptions

	err = s.data.UpdateOutpatient(existingOutpatient)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) CancelOutpatient(outpatientId int, userId int, role string) error {
	const op errors.Op = "schedules.business.CancelOutpatient"
	var errMsg errors.ErrClientMessage

	existingOutpatient, err := s.data.SelectOutpatientById(outpatientId)
	if err != nil {
		return errors.E(err, op)
	}

	if existingOutpatient.Status != schedules.StatusWaiting {
		errMsg = "Cannot cancel outpatient that is not waiting"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}

	switch role {
	case "admin":
		// just proceed to next step since admin can cancel any outpatient
	case "doctor":
		if userId != existingOutpatient.WorkSchedule.Doctor.ID {
			errMsg = "Only doctor of this outpatient work schedule can cancel this outpatient"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
		}
	case "nurse":
		if userId != existingOutpatient.WorkSchedule.Nurse.ID {
			errMsg = "Only nurse of this outpatient work schedule can cancel this outpatient"
			return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
		}

	default:
		errMsg = "You are not authorized to cancel this outpatient"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnauthorized)
	}

	existingOutpatient.Status = schedules.StatusCanceled
	err = s.data.UpdateOutpatient(existingOutpatient)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

func (s *scheduleBusiness) RemoveOutpatientById(outpatientId int) error {
	const op errors.Op = "schedules.business.RemoveOutpatientById"
	var errMsg errors.ErrClientMessage

	existingOutpatient, err := s.data.SelectOutpatientById(outpatientId)
	if err != nil {
		return errors.E(err, op)
	}

	if existingOutpatient.Status == schedules.StatusOnprogress {
		errMsg = "Cannot remove outpatient that is on progress"
		return errors.E(errors.New(string(errMsg)), op, errMsg, errors.KindUnprocessable)
	}

	err = s.data.DeleteOutpatientById(outpatientId)
	if err != nil {
		return errors.E(err, op)
	}
	return nil
}

func (s *scheduleBusiness) RemovePatientWaitingOutpatients(patientId int) error {
	const op errors.Op = "schedules.business.RemovePatientWaitingOutpatients"

	err := s.data.DeleteWaitingOutpatientsByPatientId(patientId)
	if err != nil {
		return errors.E(err, op)
	}

	return nil
}

// Private methods
func (s *scheduleBusiness) repeatEveryDay(start string, end string) ([]string, error) {
	const op errors.Op = "schedules.business.repeatEveryDay"
	const INCREMENT_DAY = 1

	dates, err := s.generateDates(start, end, INCREMENT_DAY)
	if err != nil {
		return []string{}, errors.E(err, op)
	}
	return dates, nil
}

func (s *scheduleBusiness) repeatEveryWeek(start string, end string) ([]string, error) {
	const op errors.Op = "schedules.business.repeatEveryWeek"
	const INCREMENT_DAY = 7

	dates, err := s.generateDates(start, end, INCREMENT_DAY)
	if err != nil {
		return []string{}, errors.E(err, op)
	}
	return dates, nil
}

func (s *scheduleBusiness) repeatEveryMonthSameDay(start string, end string) ([]string, error) {
	const op errors.Op = "schedules.business.repeatEveryMonthSameDay"
	const INCREMENT_DAY = 28

	dates, err := s.generateDates(start, end, INCREMENT_DAY)
	if err != nil {
		return []string{}, errors.E(err, op)
	}
	return dates, nil
}

func (s *scheduleBusiness) generateDates(start string, end string, incrementDay int) ([]string, error) {
	const op errors.Op = "schedules.business.generateDates"
	var errMessage errors.ErrClientMessage = "Invalid date format"

	var dates []string
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return dates, errors.E(err, op, errMessage, errors.KindBadRequest)
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return dates, errors.E(err, op, errMessage, errors.KindBadRequest)
	}

	// add 1 day to include the end date
	limit := endDate.AddDate(0, 0, 1)

	for d := startDate; d.Before(limit); d = d.AddDate(0, 0, incrementDay) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates, nil
}

func (s *scheduleBusiness) getUniqueSchedules(outpatients []schedules.OutpatientCore) []schedules.WorkScheduleCore {
	schedulesMap := make(map[int]schedules.WorkScheduleCore)
	for i := range outpatients {
		scheduleId := outpatients[i].WorkSchedule.ID
		schedulesMap[scheduleId] = outpatients[i].WorkSchedule
	}

	schedulesData := make([]schedules.WorkScheduleCore, len(schedulesMap))
	i := 0
	for _, schedule := range schedulesMap {
		schedulesData[i] = schedule
		i++
	}

	return schedulesData
}

func (s *scheduleBusiness) getUniqueDoctorIds(ws []schedules.WorkScheduleCore) []int {
	doctorsMap := make(map[int]bool)
	for _, s := range ws {
		doctorsMap[s.Doctor.ID] = true
	}

	doctorIds := make([]int, len(doctorsMap))
	i := 0
	for key := range doctorsMap {
		doctorIds[i] = key
		i++
	}

	return doctorIds
}

func (s *scheduleBusiness) getUniqueNurseIds(ws []schedules.WorkScheduleCore) []int {
	nursesMap := make(map[int]bool)
	for _, s := range ws {
		nursesMap[s.Nurse.ID] = true
	}

	nurseIds := make([]int, len(nursesMap))
	i := 0
	for key := range nursesMap {
		nurseIds[i] = key
		i++
	}

	return nurseIds
}

func (s *scheduleBusiness) getUniquePatientIds(outpatients []schedules.OutpatientCore) []int {
	outpatientsMap := make(map[int]bool)
	for _, o := range outpatients {
		outpatientsMap[o.Patient.ID] = true
	}

	patientIds := make([]int, len(outpatientsMap))
	i := 0
	for key := range outpatientsMap {
		patientIds[i] = key
		i++
	}

	return patientIds
}

func (s *scheduleBusiness) findDoctorsData(ids []int) (map[int]schedules.DoctorCore, error) {
	const op errors.Op = "schedules.business.findDoctorsData"

	doctorsMap := make(map[int]schedules.DoctorCore)
	doctorsData, err := s.doctorBusiness.FindDoctosrByIds(ids)
	if err != nil {
		return map[int]schedules.DoctorCore{}, errors.E(err, op)
	}

	for _, d := range doctorsData {
		doctorsMap[d.ID] = schedules.DoctorCore{
			ID:        d.ID,
			Name:      d.Name,
			Email:     d.Email,
			Phone:     d.Phone,
			Specialty: d.Speciality.Name,
			Age:       d.Age,
			Gender:    d.Gender,
			Room:      schedules.RoomCore{ID: d.Room.ID, Code: d.Room.Code, Floor: d.Room.Floor},
		}
	}

	return doctorsMap, nil
}

func (s *scheduleBusiness) findNursesData(ids []int) (map[int]schedules.NurseCore, error) {
	const op errors.Op = "schedules.business.findNursesData"

	nursesMap := make(map[int]schedules.NurseCore)
	nursesData, err := s.nurseBusiness.FindNursesByIds(ids)
	if err != nil {
		return map[int]schedules.NurseCore{}, errors.E(err, op)
	}

	for _, n := range nursesData {
		nursesMap[n.ID] = schedules.NurseCore{
			ID:     n.ID,
			Name:   n.Name,
			Email:  n.Email,
			Phone:  n.Phone,
			Age:    n.Age,
			Gender: n.Gender,
		}
	}

	return nursesMap, nil
}

func (s *scheduleBusiness) findPatientData(ids []int) (map[int]schedules.PatientCore, error) {
	const op errors.Op = "schedules.business.findPatientData"

	patientsMap := make(map[int]schedules.PatientCore)
	patientsData, err := s.patientBusiness.FindPatientsByIds(ids)
	if err != nil {
		return map[int]schedules.PatientCore{}, errors.E(err, op)
	}

	for _, p := range patientsData {
		patientsMap[p.ID] = schedules.PatientCore{
			ID:     p.ID,
			NIK:    p.NIK,
			Name:   p.Name,
			Phone:  p.Phone,
			Age:    p.Age,
			Gender: p.Gender,
		}
	}

	return patientsMap, nil
}
