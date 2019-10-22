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

func ResetPassword(nik string) {

}
func ChangePassword(nik string, password string) {

}

func GetRandomPassword(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//LOGIN SECTION

func Login(nik string, password string) (bool, interface{}, string) {
	passwordHash := getPassword(nik)
	niks := ""
	divisi := 0
	var jadwal sql.NullString
	if common.IsPasswordAndHashOk([]byte(password), passwordHash) {
		sqlLogin := `SELECT t1."NIK",t2."Divisi",t2."Jadwal" FROM "UserLogin" t1 INNER JOIN
    "ListKaryawan" t2 ON t1."NIK" = t2."NIK" WHERE t1."Password" = $1 AND t1."NIK" = $2`
		err := con.QueryRow(sqlLogin, passwordHash, nik).Scan(&niks, &divisi, &jadwal)
		if err != nil {
			return true, nil, err.Error()
		}
		tokenizer := model.LoginData{
			NIK:    niks,
			Divisi: divisi,
			Jadwal: jadwal.String,
		}
		if niks != "" && divisi != 0 {
			return false, tokenizer, "Success"
		}
		return true, nil, "UNAuthorized"

	}
	return true, nil, "Not Found!"
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
