package overtime

import (
	"absensi-server/core/action/overtime/model"
	"absensi-server/libs/database"
	deta "absensi-server/util/data"
	"database/sql"
	"fmt"
	"time"
)

var con *sql.DB

func init() {
	con = database.Connect()
}

func CreateOvertime(nik string, offset int, msg string, divisi int) (interface{}, string) {
	id := deta.RandomOvertimeID()
	ids := ""
	nowDate := time.Now()
	sqlS := `INSERT INTO "ListOvertime"("ID","Offset","DateIssued","Asignee","Message","Divisi","Status")
			VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "ID"`
	err := con.QueryRow(sqlS, id, offset, nowDate.Format("2006-01-02"), nik, msg, divisi, "WAITING").Scan(&ids)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error on inserting!"
	}
	if ids != "" {
		cover := model.CreateOvertime{
			ID:       ids,
			Assignee: nik,
			Message:  msg,
		}
		return cover, "Success!"
	}

	return nil, ""
}

func UpdateOvertime(nik string, offset int, msg string, division int, status string, id string) (interface{}, string) {
	ids := ""
	var sqlS string
	var err error
	if status != "" {
		sqlS = `UPDATE "ListOvertime" SET "Validator" = $1,"Offset" = $2,"Message"=$3,"Divisi" = $4,"Status" = $5 
			WHERE "ID" = $6 RETURNING "ID"`
		err = con.QueryRow(sqlS, nik, offset, msg, division, status, id).Scan(&ids)
	} else {
		sqlS = `UPDATE "ListOvertime" SET "Validator" = $1,"Offset" = $2,"Message"=$3,"Divisi" = $4 
			WHERE "ID" = $5 RETURNING "ID"`
		err = con.QueryRow(sqlS, nik, offset, msg, division, id).Scan(&ids)
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Some error occurred!"
	}
	if ids != "" {
		type over struct {
			ID      string
			Message string
			Status  string
		}

		ov := over{
			ID:      ids,
			Message: msg,
			Status:  status,
		}
		return ov, "Successfully!"
	}
	return nil, ""
}

func ReadListOvertime(divisi int, grant string) (interface{}, string) {
	sqlS := ""
	sqlS = `SELECT "ID","Offset","DateIssued","Asignee","Validator","Message" FROM "ListOvertime" WHERE "Divisi" = $1 AND "Status" = $2`
	rows, err := con.Query(sqlS, divisi, grant)

	if err != nil {
		fmt.Println(err.Error())
		return nil, "No return"
	}
	var j model.OvertimeList
	for rows.Next() {
		var id, date, assignee, message string
		var validator sql.NullString
		var offset int
		err = rows.Scan(&id, &offset, &date, &assignee, &validator, &message)
		if err != nil {
			fmt.Println(err.Error())
			return nil, "Error on listing"
		}
		x := model.OvertimeLister{
			ID:         id,
			Offset:     offset,
			DateIssued: date,
			Assignee:   assignee,
			Validator:  validator.String,
			Message:    message,
		}
		j.List = append(j.List, x)
	}
	if len(j.List) != 0 {
		return j, "Success"
	}
	return nil, "Error!"
}
