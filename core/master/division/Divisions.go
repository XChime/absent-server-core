package division

import (
	"absensi-server/libs/database"
	"database/sql"
	"strconv"
)

var con *sql.DB

func init() {
	con = database.Connect()
}

func showDivision() (bool, interface{}, string) {
	sqls := `SELECT "ID","NamaDivisi" FROM "ListDivisi"`
	rows, err := con.Query(sqls)
	if err != nil {
		println(err.Error())
		return true, nil, "Error!"
	}
	var dats []interface{}
	for rows.Next() {
		var id int
		var nameDivision string
		err = rows.Scan(&id, &nameDivision)
		if err != nil {
			println("ERROR ITERATING!")
		}
		dats = append(dats, ListDivision{
			ID:           id,
			NameDivision: nameDivision,
		})
	}
	if len(dats) != 0 {
		return false, dats, "Success!"
	}
	return true, nil, "No data!"
}

func createDivision(nameDivision string) (bool, string) {
	id := lastIDDivision()
	ids := 0
	sqls := `INSERT INTO "ListDivisi"("ID", "NamaDivisi") VALUES ($1,$2) RETURNING "ID"`
	err := con.QueryRow(sqls, id, nameDivision).Scan(&ids)
	if err != nil {
		println(err.Error())
		return true, "Error on inserting!"
	}
	if ids != 0 {
		return false, "Success inserting " + strconv.Itoa(ids) + " Division " + nameDivision
	}
	return true, "No data affected!"
}

func lastIDDivision() int {
	id := 0
	sqls := `SELECT "ID" FROM "ListDivisi" ORDER BY "ID" DESC LIMIT 1`
	row := con.QueryRow(sqls)
	err := row.Scan(&id)
	switch err {
	case nil:
		id = id + 1
		return id
	}
	return id
}
