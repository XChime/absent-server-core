package data

import (
	"absensi-server/libs/database"
	"database/sql"
	"fmt"
)

var con *sql.DB

func init() {
	con = database.Connect()
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

func IsNikAvailable() bool {
	return false
}
