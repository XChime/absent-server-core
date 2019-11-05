package model

type JSONSchedule struct {
	Schedule []JSONScheduleData `json:"Schedule"`
}

type JSONScheduleEmployee struct {
	Schedule []JSONScheduleDataEmployee `json:"Schedule"`
}

type JSONScheduleData struct {
	Date     string `json:"Date"`
	Shift    int    `json:"Shift"`
	Overtime string `json:"Overtime"`
	OutTime  string
	InTime   string
}

type JSONScheduleDataEmployee struct {
	Date           string `json:"Date"`
	Shift          int    `json:"Shift"`
	Overtime       string `json:"Overtime"`
	OvertimeDetail interface{}
	OutTime        string
	InTime         string
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
type OvertimeData struct {
	ID         string
	Offset     int
	DateIssued string
	Validator  string
	Message    string
	Expired    bool
}
