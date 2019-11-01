package machine

import (
	"absensi-server/core/master/machine/model"
	"absensi-server/libs/database"
	"absensi-server/util/common"
	"database/sql"
	"fmt"
)

var con *sql.DB

func init() {
	con = database.Connect()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func loginMachine(sharecode string) (interface{}, string) {
	ids := ""
	name := ""
	secret := ""
	sqlS := `SELECT "ID","Name","SecretCode" FROM "ListMachine" WHERE "ShareCode" = $1`
	err := con.QueryRow(sqlS, sharecode).Scan(&ids, &name, &secret)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "Error occurred! with sharecode " + sharecode
	}
	if ids != "" && name != "" && secret != "" {
		secret = common.HashAndSalt([]byte(secret + ids))
		datas := model.LoginMachineStr{
			IDMachine: ids,
			Name:      name,
			Secret:    secret,
		}
		return datas, "Machine FOUND! and authorized"
	}

	return nil, "ERROR Machine Login"
}

func makeShareableMachineID(idmachine string) {

}
