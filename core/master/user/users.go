package user

import (
	"absensi-server/core/master/user/model"
	"absensi-server/libs/database"
	"absensi-server/util/common"
	"absensi-server/util/data"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var con *sql.DB

func init() {
	con = database.Connect()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//Create Employee
func Create(nameCreate string, divisi string) (interface{}, string) {
	year := strconv.Itoa(common.GetYear())
	lastNik := data.CheckLastNik()
	nik, _ := strconv.Atoi(lastNik[4:])
	divisis, _ := strconv.Atoi(divisi)
	password := GetRandomPassword(8) // 8 digit random password
	if divisis < 9 {
		divisi = "0" + divisi
	}

	if nik < 9 {
		lastNik = "000" + strconv.Itoa(nik+1)
	} else if nik < 99 && nik > 9 {
		lastNik = "00" + strconv.Itoa(nik+1)
	} else if nik < 999 && nik > 99 {
		lastNik = "0" + strconv.Itoa(nik+1)
	} else {
		lastNik = strconv.Itoa(nik + 1)
	}

	newNik := fmt.Sprintf("%s%s%s", year[2:], divisi, lastNik)

	//INSERT to ListKaryawan
	createKaryawan := `INSERT INTO "ListKaryawan"("NIK","Nama","Divisi") VALUES($1,$2,$3) RETURNING "NIK"`
	NIK := ""
	err := con.QueryRow(createKaryawan, newNik, nameCreate, divisi).Scan(&NIK)
	if err != nil {
		return nil, err.Error()
	}
	createDefaultLogin := `INSERT INTO "UserLogin"("NIK","Password") VALUES($1,$2) RETURNING "Password"`
	pass := ""
	if NIK != "" {
		err := con.QueryRow(createDefaultLogin, newNik, common.HashAndSalt([]byte(password))).Scan(&pass)
		if err != nil {
			return nil, err.Error()
		}
	}
	if NIK != "" && pass != "" {
		successCreate := model.CreatedEmployee{
			NIK:             newNik,
			DefaultPassword: password,
			Message:         "Segera ganti password!",
		}
		return successCreate, "Success create karyawan!"
	}
	return nil, "ERROR"
}

func ResetAccount(nik string) (interface{}, string) {
	password := GetRandomPassword(8)
	var dbpassword sql.NullString
	sqls := `UPDATE "UserLogin" SET "Password" = $1 , "DeviceID"= NULL WHERE "NIK" = $2 RETURNING "Password"`
	err := con.QueryRow(sqls, common.HashAndSalt([]byte(password)), nik).Scan(&dbpassword)
	if err != nil {
		return nil, err.Error()
	}
	if dbpassword.Valid {
		successCreate := model.CreatedEmployee{
			NIK:             nik,
			DefaultPassword: password,
			Message:         "Account reset successfully!",
		}
		return successCreate, "Success!"
	}
	return nil, "Error!"

}
func ChangePassword(nik string, password string) (bool, string) {
	var dbpass string
	sqls := `UPDATE "UserLogin" SET "Password" = $1 WHERE "NIK" = $2 RETURNING "Password"`
	err := con.QueryRow(sqls, common.HashAndSalt([]byte(password)), nik).Scan(&dbpass)
	if err != nil {
		return false, err.Error()
	}
	ok := common.IsPasswordAndHashOk([]byte(password), dbpass)
	if ok {
		return true, "Password has been changed!"
	}

	return false, "Error!"
}

func GetRandomPassword(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

//LOGIN SECTION

func Login(nik string, password string, deviceId string) (bool, interface{}, string) {
	passwordHash := getPassword(nik)
	niks := ""
	divisi := 0
	devicehash := ""
	var deviceids sql.NullString
	var jadwal sql.NullString
	name := ""
	divisiname := ""

	if common.IsPasswordAndHashOk([]byte(password), passwordHash) {
		sqlLogin := `SELECT t1."NIK",t2."Divisi",t2."Jadwal",t1."DeviceID" ,t2."Nama",t3."NamaDivisi" FROM "UserLogin" t1 INNER JOIN
    "ListKaryawan" t2 ON t1."NIK" = t2."NIK" INNER JOIN "ListDivisi" t3 ON t2."Divisi" = t3."ID" 
	WHERE t1."Password" = $1 AND t1."NIK" = $2`

		err := con.QueryRow(sqlLogin, passwordHash, nik).Scan(&niks, &divisi, &jadwal, &deviceids, &name, &divisiname)

		if err != nil {
			return true, nil, err.Error()
		}

		if deviceids.Valid {
			if common.IsPasswordAndHashOk([]byte(deviceids.String), common.HashAndSalt([]byte(deviceId))) {
				devicehash = common.HashAndSalt([]byte(deviceids.String))
			} else {
				return true, nil, "Your account has login on another device,please call Personalia" +
					" or IT department to reset your account"
			}
		} else {
			var devid sql.NullString
			sqlUpdateDevice := `UPDATE "UserLogin" SET "DeviceID" = $1 WHERE "NIK" = $2 RETURNING "DeviceID"`
			errupdatedev := con.QueryRow(sqlUpdateDevice, deviceId, nik).Scan(&devid)
			if errupdatedev != nil {
				return true, nil, errupdatedev.Error()
			}
			devicehash = common.HashAndSalt([]byte(devid.String))
		}

		tokenizer := model.LoginData{
			NIK:        niks,
			Nama:       name,
			DivisiName: divisiname,
			Divisi:     divisi,
			Jadwal:     jadwal.String,
			DeviceHash: devicehash,
		}
		if niks != "" && divisi != 0 {
			return false, tokenizer, "Success"
		}
		return true, nil, "UNAuthorized"

	}
	return true, nil, "UNAuthorized"
}
func LoginAdmin(nik string, password string) (bool, interface{}, string) {
	passwordHash := getPassword(nik)
	niks := ""
	divisi := 0
	name := ""
	divisiname := ""
	if common.IsPasswordAndHashOk([]byte(password), passwordHash) {
		sqlLogin := `SELECT t1."NIK",t2."Divisi" ,t2."Nama",t3."NamaDivisi" FROM "UserLogin" t1 INNER JOIN
    "ListKaryawan" t2 ON t1."NIK" = t2."NIK" INNER JOIN "ListDivisi" t3 ON t2."Divisi" = t3."ID" 
	WHERE t1."Password" = $1 AND t1."NIK" = $2 AND (t2."Divisi" = 1 OR t2."Divisi" = 2)`

		err := con.QueryRow(sqlLogin, passwordHash, nik).Scan(&niks, &divisi, &name, &divisiname)

		if err != nil {
			fmt.Println(err.Error())
			return true, nil, "UnAuthorized"
		}

		tokenizer := model.AdminData{
			IsAdmin:    true,
			NIK:        niks,
			Nama:       name,
			DivisiName: divisiname,
			Divisi:     divisi,
		}
		return false, tokenizer, "Admin found"
	}
	return true, nil, ""
}

func GetEmployeeByDivision(div string) (interface{}, string) {
	sqlS := `SELECT "NIK","Nama","Divisi","NamaDivisi" FROM "ListKaryawan" t1 INNER JOIN "ListDivisi" t2 ON t1."Divisi" = t2."ID" WHERE "Divisi" = $1`
	rows, err := con.Query(sqlS, div)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Not found!"
	}
	listed := model.EmployeList{}
	for rows.Next() {
		var nik, name, namadiv string
		var div int
		err = rows.Scan(&nik, &name, &div, &namadiv)
		if err != nil {
			return nil, err.Error()
		}
		emplist := model.EmployeeListDetail{
			NIK:        nik,
			Nama:       name,
			DivisiName: namadiv,
			Divisi:     div,
		}
		listed.List = append(listed.List, emplist)
	}
	if len(listed.List) != 0 {
		return listed, "Success"
	}

	return nil, "ERROR !"
}

func ShowEmployeeProfile(nik string) (interface{}, string) {
	var niks, name, jadwal, namadivisi string
	var divisi int
	sqlS := `SELECT t1."NIK",t1."Nama",t1."Jadwal",t1."Divisi",LD."NamaDivisi" FROM "ListKaryawan" t1 
    		INNER JOIN "ListDivisi" LD on t1."Divisi" = LD."ID" WHERE "NIK" =$1`
	err := con.QueryRow(sqlS, nik).Scan(&niks, &name, &jadwal, &divisi, &namadivisi)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error on Fetching!"
	}
	if niks != "" {
		prof := model.ProfileEmployee{
			NIK:        niks,
			Nama:       name,
			Divisi:     divisi,
			NamaDivisi: namadivisi,
			Jadwal:     jadwal,
		}
		return prof, "Success!"
	}
	return nil, "Error occurred!"
}

func getPassword(nik string) string {
	var password string
	sqlStatement := `SELECT "Password" FROM "UserLogin" WHERE "NIK" = $1`
	row := con.QueryRow(sqlStatement, nik)
	switch err := row.Scan(&password); err {
	case sql.ErrNoRows:
		return password
	case nil:
		return password
	default:
		fmt.Println(err)
		return password
	}
}
