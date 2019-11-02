package model

type CreateOvertime struct {
	ID       string
	Assignee string
	Message  string
}

type OvertimeList struct {
	List []OvertimeLister
}

type OvertimeLister struct {
	ID         string
	Offset     int
	DateIssued string
	Assignee   string
	Validator  string
	Message    string
}
