package model

type AbsentEmployee struct {
	Absent []AbsentDayDetail
}

type AbsentDayDetail struct {
	Name        string
	Status      string
	Date        string
	InMachine   string
	OutMachine  string
	IN          string
	OUT         string
	Info        string
	Overtime    AbsentOvertimeDetail
	OutDetailID string
}

type AbsentOvertimeDetail struct {
	ID     string
	Offset int
}
