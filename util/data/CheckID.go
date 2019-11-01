package data

import (
	"absensi-server/libs/database"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

var con *sql.DB
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	con = database.Connect()
}

func RandomScheduleID() string {
	id := GetRandomString(10)
	if IsScheduleIDNone(id) {
		return id
	} else {
		RandomScheduleID()
	}
	return ""
}
func IsScheduleIDNone(id string) bool {
	ids := ""
	sqlS := `SELECT "ID" FROM "ListJadwal" WHERE "ID" = $1`
	row := con.QueryRow(sqlS, id)
	switch err := row.Scan(&ids); err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		fmt.Println(err)
		return true
	}
}

func CheckLastNik() string {
	sqlStatement := `SELECT "NIK" FROM "ListKaryawan" ORDER BY "NIK" DESC LIMIT 1`
	var nik string
	row := con.QueryRow(sqlStatement)
	switch err := row.Scan(&nik); err {
	case sql.ErrNoRows:
		return nik
	case nil:
		return nik
	default:
		fmt.Println(err)
		return nik
	}
}

func IsNIKAdmin(nik string) bool {
	div := 0
	sqlStatement := `SELECT "Divisi" FROM "ListKaryawan" WHERE "NIK" = $1`
	row := con.QueryRow(sqlStatement, nik)
	switch err := row.Scan(&div); err {
	case sql.ErrNoRows:
		return false
	case nil:
		if div == 1 || div == 2 {
			return true
		}
		return false
	default:
		fmt.Println(err.Error())
		return false
	}
}

func GetRandomString(n int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
