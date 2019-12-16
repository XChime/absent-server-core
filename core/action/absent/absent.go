/*Copyright 2019
Davin Alfarizky Putra Basudewa <dbasudewa@gmail.com,moshi2_davin@dvnlabs.ml>
absensi-server corw written in GO
*/
package absent

import (
	medal "absensi-server/core/action/absent/model"
	"absensi-server/core/action/schedule/model"
	"absensi-server/libs/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"strconv"
	"time"
)

var con *sql.DB

func init() {
	con = database.Connect()
}

func RequestAbsent(clientID string, machineID string, deviceID string) (bool, string) {
	isSuccess := false
	message := ""
	var dates = time.Now()
	requesType := RequestType(clientID, dates)
	if requesType == "IN" {
		println("IN")
		isSuccess, message = InRequest(clientID, machineID, deviceID, dates)
	} else if requesType == "OUT" {
		println("OUT")
		isSuccess, message = OutRequest(clientID, machineID, deviceID, dates)
	} else if requesType == "INS" {
		println("BackToOffice REQUEST")
		isSuccess, message = INSRequest(clientID, machineID, deviceID, dates)
	} else if requesType == "OVERTIME" {
		println("Overtime REQUEST")
		isSuccess, message = OvertimeAbsentRequest(clientID, machineID, deviceID, dates)
	} else if requesType == "" {
		return false, "Discipline may lead to your success!"
	}

	return isSuccess, message
}

//This function will be called every time on making request using login machine
func RequestType(nik string, datetime time.Time) string {
	//Return IN if absensi not found today
	//Return OUT if status LATE or IN
	//Return INS if status OUT and have OutDetail not have Overtime
	//Return Overtime if status OUT and have Overtime
	var status, in string
	var out, overtime sql.NullString
	date := datetime.Format("2006-01-02")
	sqlS := `SELECT "Status","OutDetailsID","InMachineID","Overtime" FROM "ListAbsensi" WHERE "NIK" = $1 AND "Date" = $2`
	err := con.QueryRow(sqlS, nik, date).Scan(&status, &out, &in, &overtime)
	if err != nil {
		fmt.Println(err.Error())
		return "IN"
	}
	if status == "" {
		return "IN"
	} else if (status == "IN" || status == "LATE") && in != "" {
		return "OUT"
	} else if status == "OUT" && out.String != "" && overtime.String == "" {
		return "INS"
	} else if (status == "IN" || status == "INF" || status == "LATE") && overtime.String != "" {
		return "OVERTIME"
	}
	return ""
}
func OvertimeAbsentRequest(clientID, machineID, deviceID string, dates time.Time) (bool, string) {
	nowTime := dates.Format("3:04 PM")
	date := dates.Format("2006-01-02")
	jadwal := ""
	offset := 0
	sqlS := `SELECT LO."Offset",LJ."ID" FROM "ListAbsensi" t1 
    		INNER JOIN "ListKaryawan" LK on t1."NIK" = LK."NIK"
    		INNER JOIN "ListJadwal" LJ on LK."Jadwal" = LJ."ID"
    		INNER JOIN "ListOvertime" LO on t1."Overtime" = LO."ID"
    		INNER JOIN "UserLogin" UL on LK."NIK" = UL."NIK"
    		WHERE t1."NIK" = $1 AND LO."Status" = 'GRANTED' AND UL."DeviceID" = $2`
	err := con.QueryRow(sqlS, clientID, deviceID).Scan(&offset, &jadwal)
	if err != nil {
		fmt.Println(err.Error())
		return false, "INIT FAILED!"
	}
	if jadwal != "" && offset != 0 {
		isOver := false
		var tmout, now time.Time
		schedules := ReadSchedule(jadwal)
		if schedules.Schedule != nil {
			for i := 0; i < len(schedules.Schedule); i++ {
				da := schedules.Schedule[i]
				tmout, _ = time.Parse("3:04 PM", da.OutTime)
				now, _ = time.Parse("3:04 PM", nowTime)
				if now.After(tmout) {
					isOver = true
				}
			}
			if isOver {
				//Checking difference between time out to time request.If over than 1 hours,
				//Request will be granted.
				diffOffset := now.Sub(tmout).Hours()
				if diffOffset >= 1 {
					niks := ""
					sqlE := `UPDATE "ListAbsensi" SET "Status" = $1,"OutMachineID" = $2,"OUT" = $3 
							WHERE "NIK" = $4 AND "Date" = $5 RETURNING "NIK"`
					errs := con.QueryRow(sqlE, "OUT", machineID, nowTime, clientID, date).Scan(&niks)
					if errs != nil {
						fmt.Println(errs.Error())
						return false, "Error Updating!"
					}
					if niks != "" {
						return true, "お疲れさまでした"
					}
				}
			} else {
				return false, "Do you want to cheating?!"
			}
		}
	}

	return false, "NO DATA!"
}

func INSRequest(clientID, machineID, deviceID string, dates time.Time) (bool, string) {
	nowTime := dates.Format("3:04 PM")
	date := dates.Format("2006-01-02")
	jadwal := ""
	sqlS := `SELECT LJ."ID" FROM "ListAbsensi" t1 
    		INNER JOIN "ListKaryawan" t2 ON t1."NIK" = t2."NIK" 
    		INNER JOIN "ListJadwal" LJ on t2."Jadwal" = LJ."ID"
			INNER JOIN "OutDetails" t3 ON t1."OutDetailsID" = t3."ID"
			INNER JOIN "UserLogin" t4 ON t2."NIK" = t4."NIK"
			WHERE t1."NIK" = $1 AND t3."Status" = 'GRANTED' AND t4."DeviceID" = $2`
	err := con.QueryRow(sqlS, clientID, deviceID).Scan(&jadwal)
	if err != nil {
		fmt.Println(err.Error())
		return false, "Error INS"
	}
	if jadwal != "" {
		isBetween := false
		schedules := ReadSchedule(jadwal)
		if schedules.Schedule != nil {
			for i := 0; i < len(schedules.Schedule); i++ {
				da := schedules.Schedule[i]
				tmin, _ := time.Parse("3:04 PM", da.InTime)
				tmout, _ := time.Parse("3:04 PM", da.OutTime)
				tm2, _ := time.Parse("3:04 PM", nowTime)
				if tm2.After(tmin) && tm2.Before(tmout) {
					isBetween = true
				}
			}
			if isBetween {
				niks := ""
				sqlU := `UPDATE "ListAbsensi" SET "Status" = $1,"InMachineID" = $2,"Info" = $3 
						WHERE "NIK" = $4 AND "Date" = $5 RETURNING "NIK"`
				errs := con.QueryRow(sqlU, "INF", machineID, "BACK TO OFFICE", clientID, date).Scan(&niks)
				if errs != nil {
					fmt.Println(errs.Error())
					return false, "Error Back to office"
				}
				if niks != "" {
					return true, "Success Back to Office!"
				}
			} else {
				return false, "Go back yourself!"
			}
		}
	}

	return false, "Error No Data!"
}

func OutRequest(clientID, machineID, deviceID string, dates time.Time) (bool, string) {
	var idout, assignee, message, jadwal string
	date := dates.Format("2006-01-02")
	nowTime := dates.Format("3:04 PM")
	sqlS := `SELECT t1."OutDetailsID",t2."Assignee",t2."Message",LK."Jadwal" FROM "ListAbsensi"t1 
			INNER JOIN "OutDetails" t2 ON t1."OutDetailsID" = t2."ID"
			INNER JOIN "UserLogin" UL on t1."NIK" = UL."NIK"
			INNER JOIN "ListKaryawan" LK ON LK."NIK" = t1."NIK"
			WHERE t2."Status" = 'GRANTED' AND t1."NIK" = $1 AND UL."DeviceID" = $2`
	err := con.QueryRow(sqlS, clientID, deviceID).Scan(&idout, &assignee, &message, &jadwal)
	if err != nil {
		fmt.Println(err.Error())
		return false, "ERROR"
	}
	if idout != "" && assignee != "" && message != "" && jadwal != "" {
		schedules := ReadSchedule(jadwal)
		isBack := true
		if schedules.Schedule != nil {
			for i := 0; i < len(schedules.Schedule); i++ {
				da := schedules.Schedule[i]
				tmin, _ := time.Parse("3:04 PM", da.InTime)
				tmout, _ := time.Parse("3:04 PM", da.OutTime)
				tm2, _ := time.Parse("3:04 PM", nowTime)
				if tm2.After(tmin) && tm2.Before(tmout) {
					isBack = false
				}
			}
			if isBack {
				return false, "戻ってがいます"
			}
			status := ""
			sqlI := `UPDATE "ListAbsensi" SET "Status" = 'OUT',"OutMachineID" = $1,"OUT" = $2 
					WHERE "NIK" = $3 AND "Date" = $4 RETURNING "Status"`
			errs := con.QueryRow(sqlI, machineID, nowTime, clientID, date).Scan(&status)
			if errs != nil {
				fmt.Println(errs.Error())
				return false, "Error!"
			}
			if status == "OUT" {
				return true, "Success OUT Request!"
			}
		}
	}
	return false, ""
}

func InRequest(clientID string, machineID string, deviceID string, dates time.Time) (bool, string) {
	ids := ""
	jadwal := ""
	sqlS := `SELECT t1."NIK",t2."Jadwal" FROM "UserLogin" t1 INNER JOIN "ListKaryawan" t2 ON t1."NIK" = t2."NIK"
			WHERE t1."NIK" = $1 AND t1."DeviceID" = $2`
	err := con.QueryRow(sqlS, clientID, deviceID).Scan(&ids, &jadwal)
	if err != nil {
		fmt.Println(err.Error())
	}

	if ids == "" && jadwal == "" {
		return false, "Device ID or Client ID not authorized,please reset account!"
	}
	schedules := ReadSchedule(jadwal)
	if schedules.Schedule != nil {
		isLate := true
		niks := ""
		date := dates.Format("2006-01-02")
		inTime := dates.Format("3:04 PM")
		NotFoundScheduleCount := 0
		for i := 0; i < len(schedules.Schedule); i++ {
			da := schedules.Schedule[i]
			if da.Date == date {
				tmin, _ := time.Parse("3:04 PM", da.InTime)
				tmout, _ := time.Parse("3:04 PM", da.OutTime)
				tm2, _ := time.Parse("3:04 PM", inTime)
				if tm2.After(tmin) && tm2.Before(tmout) {
					isLate = true
				} else {
					isLate = false
				}
			} else {
				NotFoundScheduleCount += 1
			}
		}
		if NotFoundScheduleCount == len(schedules.Schedule) {
			return false, "Schedule Not Found!"
		}
		status := "IN"
		if isLate {
			status = "LATE"
		}
		sqlSs := `INSERT INTO "ListAbsensi"("NIK","Date","Status","InMachineID","IN") VALUES ($1,$2,$3,$4,$5) RETURNING "NIK"`
		errs := con.QueryRow(sqlSs, clientID, date, status, machineID, inTime).Scan(&niks)
		if errs != nil {
			fmt.Println(errs.Error())
			return false, "ERROR!"
		}

		if niks == "" {
			return false, "ERROR"
		}
		return true, "Success IN Request!"
	}
	return false, "ERROR!"
}

func ReadEmployeeAbsent(nik string) (interface{}, string) {
	dates := time.Now()
	monthNow := dates.Month()
	var status, inmchid, name string
	var outmchid, info, overtimeid, outdetailid, offset sql.NullString
	var date, intime time.Time
	var outtime pq.NullTime
	sqlS := `SELECT t1."Status",t1."Date",t1."InMachineID",t1."OutMachineID",t1."IN",t1."OUT"
			,t1."Info",LK."Nama",LO."ID",LO."Offset",OD."ID"
			FROM "ListAbsensi" t1 
    		INNER JOIN "ListKaryawan" LK on t1."NIK" = LK."NIK"
    		LEFT JOIN "ListOvertime" LO on t1."Overtime" = LO."ID"
    		LEFT JOIN "OutDetails" OD on t1."OutDetailsID" = OD."ID"
    		WHERE t1."NIK" = $1 AND date_part('MONTH', t1."Date") = $2 ORDER BY t1."Date" DESC`
	rows, err := con.Query(sqlS, nik, monthNow)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "ERROR!"
	}
	var k medal.AbsentEmployee
	for rows.Next() {
		err = rows.Scan(&status, &date, &inmchid, &outmchid, &intime, &outtime, &info, &name, &overtimeid, &offset, &outdetailid)
		if err != nil {
			fmt.Println(err.Error())
			return nil, "Error Increment!"
		}
		offsetI, _ := strconv.Atoi(offset.String)
		dateF := date.Format("2006-01-02")
		inF := intime.Format("3:04 PM")
		outF := ""
		if outtime.Valid {
			outF = outtime.Time.Format("3:04 PM")
		}
		x := medal.AbsentDayDetail{
			Name:       name,
			Status:     status,
			Date:       dateF,
			InMachine:  inmchid,
			OutMachine: outmchid.String,
			IN:         inF,
			OUT:        outF,
			Info:       info.String,
			Overtime: medal.AbsentOvertimeDetail{
				ID:     overtimeid.String,
				Offset: offsetI,
			},
			OutDetailID: outdetailid.String,
		}
		k.Absent = append(k.Absent, x)
	}
	if k.Absent != nil {
		return k, "Success Fetch Absent!"
	}

	return nil, "NO DATA"
}

func ReadAbsentByDay(dates string) (interface{}, string) {
	var status, inmchid, name string
	var outmchid, info, overtimeid, outdetailid, offset sql.NullString
	var date, intime time.Time
	var outtime pq.NullTime

	sqlS := `SELECT t1."Status",t1."Date",t1."InMachineID",t1."OutMachineID",t1."IN",t1."OUT"
			,t1."Info",LK."Nama",LO."ID",LO."Offset",OD."ID"
			FROM "ListAbsensi" t1 
    		INNER JOIN "ListKaryawan" LK on t1."NIK" = LK."NIK"
    		LEFT JOIN "ListOvertime" LO on t1."Overtime" = LO."ID"
    		LEFT JOIN "OutDetails" OD on t1."OutDetailsID" = OD."ID"
    		WHERE t1."Date" = $1 ORDER BY t1."NIK"`
	rows, err := con.Query(sqlS, dates)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "NO DATA"
	}
	var k medal.AbsentEmployee
	for rows.Next() {
		err = rows.Scan(&status, &date, &inmchid, &outmchid, &intime, &outtime, &info, &name, &overtimeid, &offset, &outdetailid)
		if err != nil {
			fmt.Println(err.Error())
			return nil, "Error Increment!"
		}
		offsetI, _ := strconv.Atoi(offset.String)
		dateF := date.Format("2006-01-02")
		inF := intime.Format("3:04 PM")
		outF := ""
		if outtime.Valid {
			outF = outtime.Time.Format("3:04 PM")
		}
		x := medal.AbsentDayDetail{
			Name:       name,
			Status:     status,
			Date:       dateF,
			InMachine:  inmchid,
			OutMachine: outmchid.String,
			IN:         inF,
			OUT:        outF,
			Info:       info.String,
			Overtime: medal.AbsentOvertimeDetail{
				ID:     overtimeid.String,
				Offset: offsetI,
			},
			OutDetailID: outdetailid.String,
		}
		k.Absent = append(k.Absent, x)
	}
	if k.Absent != nil {
		return k, "Success Fetch Absent!"
	}
	return nil, "NO DATA!"
}

//Other Section (NEED TO MAKE PUBLIC/GENERAL to avoid multiple function

func ReadSchedule(id string) model.JSONSchedule {
	var data string
	sqlS := `SELECT "Data" FROM "ListJadwal" WHERE "ID" = $1`
	err := con.QueryRow(sqlS, id).Scan(&data)
	if err != nil {
		fmt.Println(err.Error())
		return model.JSONSchedule{}
	}

	bytes := []byte(data)
	var j model.JSONSchedule
	errs := json.Unmarshal(bytes, &j)
	if errs != nil {
		return model.JSONSchedule{}
	}

	if len(j.Schedule) != 0 {
		for i := 0; i < len(j.Schedule); i++ {
			in, out := readShiftTime(j.Schedule[i].Shift)
			j.Schedule[i].InTime = in.Format("3:04 PM")
			j.Schedule[i].OutTime = out.Format("3:04 PM")
		}
		return j
	}

	return j
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
