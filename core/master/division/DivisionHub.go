package division

import (
	"absensi-server/core/master/user/model"
	"absensi-server/libs"
	"absensi-server/util"
	deta "absensi-server/util/data"
	"encoding/json"
)

type returnData struct {
	Error   bool
	Message string
	Data    interface{}
}
type returnNoData struct {
	Error   bool
	Message string
}

func ShowDivisionHub() interface{} {
	errs := true
	msgs := ""
	var data interface{}
	isErr, datas, msgerror := showDivision()
	if !isErr && datas != nil {
		errs = false
		data = datas
	}
	msgs = msgerror
	return returnData{
		Error:   errs,
		Message: msgs,
		Data:    data,
	}
}

func CreateDivisionHub(token string, nameDivisions string) interface{} {
	errs := true
	msgs := ""
	isErr, msgss := util.CheckParameter([]util.Parameters{{"token", token}, {"nameDivisions", nameDivisions}})
	if isErr {
		msgs = msgss
	} else {
		isOk, dat, msges := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.AdminData
			_ = json.Unmarshal(jsonString, &loginData)
			isADMINDB := deta.IsNIKAdmin(loginData.NIK)
			if isADMINDB && loginData.IsAdmin == true {
				isErr, msgerror := createDivision(nameDivisions)
				if !isErr {
					errs = false
				}
				msgs = msgerror
			} else if isADMINDB {
				errs = true
				msgs = "Required administrator token"
			} else {
				errs = true
				msgs = "Required Personalia / IT Department access"
			}
		} else {
			errs = true
			msgs = msges
		}
	}
	return returnNoData{
		Error:   errs,
		Message: msgs,
	}
}
