package response

import "github.com/final-project-alterra/hospital-management-system-api/features/schedules"

type WorkScheduleResponse struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	TotalWaiting int    `json:"totalWaiting"`

	Doctor struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Speciality string `json:"speciality"`
		BirthDate  string `json:"birthDate"`
		Gender     string `json:"gender"`

		Room struct {
			ID    int    `json:"id"`
			Code  string `json:"code"`
			Floor string `json:"floor"`
		} `json:"room"`
	} `json:"doctor"`

	Nurse struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	} `json:"nurse"`
}

type DoctorWorkScheduleResponse struct {
	ID        int    `json:"id"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`

	Nurse struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	} `json:"nurse"`
}

type NurseWorkScheduleResponse struct {
	ID        int    `json:"id"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`

	Doctor struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Speciality string `json:"speciality"`
		BirthDate  string `json:"birthDate"`
		Gender     string `json:"gender"`

		Room struct {
			ID    int    `json:"id"`
			Code  string `json:"code"`
			Floor string `json:"floor"`
		} `json:"room"`
	} `json:"doctor"`
}

func WorkSchedule(w schedules.WorkScheduleCore) WorkScheduleResponse {
	resp := WorkScheduleResponse{}

	resp.ID = w.ID
	resp.Date = w.Date
	resp.StartTime = w.StartTime
	resp.EndTime = w.EndTime
	resp.TotalWaiting = w.TotalWaiting

	resp.Doctor.ID = w.Doctor.ID
	resp.Doctor.Name = w.Doctor.Name
	resp.Doctor.Email = w.Doctor.Email
	resp.Doctor.Phone = w.Doctor.Phone
	resp.Doctor.Speciality = w.Doctor.Specialty
	resp.Doctor.BirthDate = w.Doctor.BirthDate
	resp.Doctor.Gender = w.Doctor.Gender

	resp.Doctor.Room.ID = w.Doctor.Room.ID
	resp.Doctor.Room.Code = w.Doctor.Room.Code
	resp.Doctor.Room.Floor = w.Doctor.Room.Floor

	resp.Nurse.ID = w.Nurse.ID
	resp.Nurse.Name = w.Nurse.Name
	resp.Nurse.Email = w.Nurse.Email
	resp.Nurse.Phone = w.Nurse.Phone
	resp.Nurse.Age = w.Nurse.Age
	resp.Nurse.Gender = w.Nurse.Gender

	return resp
}

func DoctorSchedule(w schedules.WorkScheduleCore) DoctorWorkScheduleResponse {
	resp := DoctorWorkScheduleResponse{}

	resp.ID = w.ID
	resp.Date = w.Date
	resp.StartTime = w.StartTime
	resp.EndTime = w.EndTime

	resp.Nurse.ID = w.Nurse.ID
	resp.Nurse.Name = w.Nurse.Name
	resp.Nurse.Email = w.Nurse.Email
	resp.Nurse.Phone = w.Nurse.Phone
	resp.Nurse.Age = w.Nurse.Age
	resp.Nurse.Gender = w.Nurse.Gender

	return resp
}

func NurseSchedule(w schedules.WorkScheduleCore) NurseWorkScheduleResponse {
	resp := NurseWorkScheduleResponse{}

	resp.ID = w.ID
	resp.Date = w.Date
	resp.StartTime = w.StartTime
	resp.EndTime = w.EndTime

	resp.Doctor.ID = w.Doctor.ID
	resp.Doctor.Name = w.Doctor.Name
	resp.Doctor.Email = w.Doctor.Email
	resp.Doctor.Phone = w.Doctor.Phone
	resp.Doctor.Speciality = w.Doctor.Specialty
	resp.Doctor.BirthDate = w.Doctor.BirthDate
	resp.Doctor.Gender = w.Doctor.Gender

	resp.Doctor.Room.ID = w.Doctor.Room.ID
	resp.Doctor.Room.Code = w.Doctor.Room.Code
	resp.Doctor.Room.Floor = w.Doctor.Room.Floor

	return resp
}

func ListWorkSchedule(w []schedules.WorkScheduleCore) []WorkScheduleResponse {
	result := make([]WorkScheduleResponse, len(w))
	for i := range w {
		result[i] = WorkSchedule(w[i])
	}
	return result
}

func ListDoctorSchedule(w []schedules.WorkScheduleCore) []DoctorWorkScheduleResponse {
	result := make([]DoctorWorkScheduleResponse, len(w))
	for i := range w {
		result[i] = DoctorSchedule(w[i])
	}
	return result
}

func ListNurseSchedule(w []schedules.WorkScheduleCore) []NurseWorkScheduleResponse {
	result := make([]NurseWorkScheduleResponse, len(w))
	for i := range w {
		result[i] = NurseSchedule(w[i])
	}
	return result
}
