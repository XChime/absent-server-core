package model

type JSONSchedule struct {
	Schedule []JSONScheduleData `json:"Schedule"`
}
type JSONScheduleData struct {
	Date     string `json:"Date"`
	Shift    int    `json:"Shift"`
	Overtime string `json:"Overtime"`
	OutTime  string
	InTime   string
}
type ScheduleDivision struct {
	List []ScheduleList
}
type ScheduleList struct {
	ID        string
	Title     string
	Division  int
	Assignee  string
	Validator string
	Schedule  []JSONScheduleData `json:"Schedule"`
}
type UpdateSchedule struct {
	ID        string
	Message   string
	Validator string
}
type CreateSchedule struct {
	ID       string
	Message  string
	Assignee string
	Division int
}
type AttachSchedule struct {
	NIK     string
	Jadwal  string
	Message string
}
