package absent

import (
	mdl "absensi-server/core/master/machine/model"
	"absensi-server/core/master/user/model"
	"absensi-server/libs"
	deta "absensi-server/util/data"
	"encoding/json"
)

func RequestAbsentHub(client string, machine string, deviceid string) interface{} {
	errs := true
	message := ""
	type absention struct {
		Error   bool
		Message string
	}

	if client != "" && machine != "" && deviceid != "" {
		isOkClient, datclient, msgs := libs.VerifyTokenGeneric(client, model.LoginData{})
		isOkMachine, datmachine, msgss := libs.VerifyTokenGeneric(machine, mdl.Access{})
		if isOkClient && isOkMachine {
			jsonString := deta.MustMarshal(datclient)
			var loginData model.LoginData
			_ = json.Unmarshal(jsonString, &loginData)
			jsonStrings := deta.MustMarshal(datmachine)
			var machineAccess mdl.Access
			_ = json.Unmarshal(jsonStrings, &machineAccess)
			reqOK, msgg := RequestAbsent(loginData.NIK, machineAccess.ID, deviceid)
			if reqOK {
				errs = false
			}
			message = msgg
		} else {
			if msgss != "" {
				message = msgss
			} else {
				message = msgs
			}
		}
	} else {
		message = "One or more field empty!"
	}

	abs := absention{
		Error:   errs,
		Message: message,
	}
	return abs
}

func ReadAbsentByDaysHub(date string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type employ struct {
		Error   bool
		Message string
		Data    interface{}
	}
	if date != "" {
		dat, msg := ReadAbsentByDay(date)
		if dat != nil {
			errs = false
			data = dat
		}
		message = msg
	}

	emp := employ{
		Error:   errs,
		Message: message,
		Data:    data,
	}

	return emp
}

func ReadAbsentByEmployeeHub(nik string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type employ struct {
		Error   bool
		Message string
		Data    interface{}
	}
	if nik != "" {
		dat, msg := ReadEmployeeAbsent(nik)
		if dat != nil {
			errs = false
			data = dat
		}
		message = msg
	}

	emp := employ{
		Error:   errs,
		Message: message,
		Data:    data,
	}

	return emp
}
