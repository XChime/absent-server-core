package schedule

import (
	"absensi-server/core/action/schedule/model"
	"absensi-server/libs/database"
	deta "absensi-server/util/data"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var con *sql.DB

func init() {
	con = database.Connect()
}
func CreateSchedule(assignee string, divisi string, mesg string, datas string) (bool, string, interface{}) {
	div, _ := strconv.Atoi(divisi)
	id := deta.RandomScheduleID()
	ids := ""
	sqlS := `INSERT INTO "ListJadwal"("ID","Divisi","Message","Data","Status","Assignee")
			VALUES ($1,$2,$3,$4,$5,$6) RETURNING "ID"`
	err := con.QueryRow(sqlS, id, div, mesg, datas, "WAITING", assignee).Scan(&ids)
	if err != nil {
		fmt.Println(err.Error())
		return false, "Error occurred", nil
	}
	if ids != "" {
		sch := model.CreateSchedule{
			ID:       ids,
			Message:  mesg,
			Assignee: assignee,
			Division: div,
		}
		return true, "Success create new Schedule", sch
	}
	return false, "Error!", nil
}

func ScheduleAttach(idsch string, nik string) (interface{}, string) {
	ids := ""
	niks := ""
	sqlS := `UPDATE "ListKaryawan" SET "Jadwal" = $1 WHERE "NIK" = $2 RETURNING "NIK","Jadwal"`
	err := con.QueryRow(sqlS, idsch, nik).Scan(&ids, &niks)

	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error processing"
	}

	if niks != "" && ids != "" {
		sch := model.AttachSchedule{
			NIK:     niks,
			Jadwal:  ids,
			Message: "Attached!",
		}
		return sch, "Success attaching employee!"
	}
	return nil, "ERROR!"
}

func ScheduleUpdate(idsch string, validator string, msgs string, datas string, grant string) (interface{}, string) {
	ids := ""
	var sqlS string
	var err error
	if grant == "GRANTED" {
		sqlS = `UPDATE "ListJadwal" SET "Validator" = $1 ,"Message" =$2,"Data"=$3,"Status" = $4 WHERE "ID" =  $5 RETURNING "ID"`
		err = con.QueryRow(sqlS, validator, msgs, datas, grant, idsch).Scan(&ids)
	} else {
		sqlS = `UPDATE "ListJadwal" SET "Validator" = $1 ,"Message" =$2,"Data"=$3 WHERE "ID" = $4 RETURNING "ID"`
		err = con.QueryRow(sqlS, validator, msgs, datas, idsch).Scan(&ids)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error processing"
	}
	if ids != "" {
		sch := model.UpdateSchedule{
			ID:        ids,
			Message:   msgs,
			Validator: validator,
		}

		return sch, "Success Updated Schedule!"
	}

	return nil, "ERROR"
}

func ReadEmployeeSchedule(nik string) (interface{}, string) {
	ids := ""
	message := ""
	namadiv := ""
	var data string
	sqlS := `SELECT t1."ID",t1."Message",LD."NamaDivisi",t1."Data" FROM "ListJadwal" t1
    INNER JOIN "ListDivisi" LD on t1."Divisi" = LD."ID" INNER JOIN "ListKaryawan" LK on t1."ID" = LK."Jadwal"
	WHERE t1."Status"= 'GRANTED' AND LK."NIK" = $1`
	err := con.QueryRow(sqlS, nik).Scan(&ids, &message, &namadiv, &data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "No Schedule attached"
	}

	bytes := []byte(data)
	var j model.JSONScheduleEmployee
	errs := json.Unmarshal(bytes, &j)
	if errs != nil {
		return nil, errs.Error()
	}

	if len(j.Schedule) != 0 {
		for i := 0; i < len(j.Schedule); i++ {
			overtime := ReadOvertime(j.Schedule[i].Overtime, j.Schedule[i].Date)
			j.Schedule[i].OvertimeDetail = overtime
			in, out := readShiftTime(j.Schedule[i].Shift)
			j.Schedule[i].InTime = in.Format("3:04 PM")
			j.Schedule[i].OutTime = out.Format("3:04 PM")
		}
		return j, message
	}
	return nil, ""
}

func ScheduleById(id string) (interface{}, string) {
	message := ""
	data := ""
	sqlS := `SELECT "Message","Data" FROM "ListJadwal" WHERE "ID" = $1 AND "Status" = 'GRANTED'`
	err := con.QueryRow(sqlS, id).Scan(&message, &data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error ID not found"
	}
	bytes := []byte(data)
	var j model.JSONSchedule
	errs := json.Unmarshal(bytes, &j)
	if errs != nil {
		return nil, errs.Error()
	}

	if len(j.Schedule) != 0 {
		for i := 0; i < len(j.Schedule); i++ {
			in, out := readShiftTime(j.Schedule[i].Shift)
			j.Schedule[i].InTime = in.Format("3:04 PM")
			j.Schedule[i].OutTime = out.Format("3:04 PM")
		}
		return j, message
	}

	return nil, "ID Not found"
}
func ScheduleDivison(div int, grant string) (interface{}, string) {
	var ids, divisi, message string
	var assginee, validator sql.NullString
	sqlS := `SELECT "ID","Divisi","Assignee","Validator","Message" FROM "ListJadwal"
			WHERE "Divisi" = $1 AND "Status" = $2`
	rows, err := con.Query(sqlS, div, grant)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err.Error()
	}
	var k model.ScheduleDivision
	for rows.Next() {
		err = rows.Scan(&ids, &divisi, &assginee, &validator, &message)
		var j model.ScheduleList
		j.Assignee = assginee.String
		j.Validator = validator.String
		divv, _ := strconv.Atoi(divisi)
		j.Division = divv
		j.ID = ids
		j.Title = message
		k.List = append(k.List, j)
	}

	if len(k.List) != 0 {
		return k, "Success"
	}
	return nil, "No List"
}
func ReadOvertime(overtime string, date string) interface{} {
	var offset int
	var validator, message string
	var dates time.Time
	sqlS := `SELECT "Offset","DateIssued","Validator","Message" FROM "ListOvertime" WHERE "ID" = $1 AND "Status" = 'GRANTED'`
	err := con.QueryRow(sqlS, overtime).Scan(&offset, &dates, &validator, &message)
	if err != nil {
		fmt.Println(overtime + " -> " + err.Error())
		return nil
	}
	dateE := dates.Format("2006-01-02")
	expired := true
	if dateE == date {
		expired = false
	}
	over := model.OvertimeData{
		ID:         overtime,
		Offset:     offset,
		DateIssued: dateE,
		Validator:  validator,
		Message:    message,
		Expired:    expired,
	}
	return over
}

func readShiftTime(shift int) (time.Time, time.Time) {
	var in time.Time
	var out time.Time
	sqlS := `SELECT "InTime","OutTime" FROM "ListTime" WHERE "Shift" = $1`
	err := con.QueryRow(sqlS, shift).Scan(&in, &out)
	if err != nil {
		return time.Time{}, time.Time{}

	}
	return in, out
}
