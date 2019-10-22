/*File untuk return connection database dari postgres
dengan konsep OOP*/

package database

import (
	"absensi-server/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

//Connect ke database server
func Connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if db != nil {
		err = db.Ping()
	}
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
		return nil
	}
	return db
}
