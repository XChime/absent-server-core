package overtime

import (
	"absensi-server/core/master/user/model"
	"absensi-server/libs"
	deta "absensi-server/util/data"
	"encoding/json"
)

func CreateOvertimeHub(token string, offset int, msg string, divisi int) interface{} {
	errs := true
	message := ""
	var data interface{}
	type createover struct {
		Error    bool
		Message  string
		Overtime interface{}
	}

	if token != "" && msg != "" && divisi != 0 {
		isOk, dat, msgs := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.AdminData
			_ = json.Unmarshal(jsonString, &loginData)
			if loginData.IsAdmin == true {
				dat, msg := CreateOvertime(loginData.NIK, offset, msg, divisi)
				if dat != nil {
					data = dat
					errs = false
					message = msg
				} else {
					message = msg
				}
			} else {
				message = "Unauthorized"
			}

		} else {
			message = msgs
		}

	} else {
		message = "One or more field no input!"
	}

	over := createover{
		Error:    errs,
		Message:  message,
		Overtime: data,
	}
	return over
}

func UpdateOvertimeHub(token string, id string, offset int, msg string, division int, status string) interface{} {
	type overtime struct {
		Error    bool
		Message  string
		Overtime interface{}
	}
	errs := true
	message := ""
	var data interface{}

	if token != "" && division != 0 && id != "" {
		isOk, dat, msgs := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.AdminData
			_ = json.Unmarshal(jsonString, &loginData)
			if loginData.IsAdmin == true {
				dat, msg := UpdateOvertime(loginData.NIK, offset, msg, division, status, id)
				if dat != nil {
					data = dat
					errs = false
					message = msg
				} else {
					message = msg
				}
			} else {
				message = "Unauthorized"
			}
		} else {
			message = msgs
		}
	} else {
		message = "One or more field not entered!"
	}

	over := overtime{
		Error:    errs,
		Message:  message,
		Overtime: data,
	}
	return over
}

func ReadOvertimeListHub(divisi int, grant string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type overtime struct {
		Error   bool
		Message string
		Data    interface{}
	}

	if divisi != 0 {
		dat, msg := ReadListOvertime(divisi, grant)
		if dat != nil {
			errs = false
			data = dat
		}
		message = msg
	} else {
		message = "One or more field not found!"
	}

	over := overtime{
		Error:   errs,
		Message: message,
		Data:    data,
	}
	return over
}
