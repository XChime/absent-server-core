package machine

import (
	"absensi-server/core/master/machine/model"
	"absensi-server/libs"
	deta "absensi-server/util/data"
)

//function untuk login pertama mesin absen
func LoginMachineHub(sharecode string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type login struct {
		Error   bool
		Message string
		Machine interface{}
	}

	if sharecode != "" {
		dat, msg := loginMachine(sharecode)
		if dat != nil {
			errs = false
			message = msg
			data = dat
		} else {
			errs = true
			message = msg
		}
	} else {
		errs = true
		message = "sharecode empty!"
	}

	returnedData := login{
		Error:   errs,
		Message: message,
		Machine: data,
	}
	return returnedData
}
func RequestMachineAccessHub(id string, secret string) interface{} {
	errs := true
	message := ""
	var data string
	type request struct {
		Error   bool
		Message string
		Token   string
	}
	if id != "" && secret != "" {
		if deta.IsMachineAccessOK(id, secret) {

			acs := model.Access{
				ID:     id,
				Secret: secret,
			}
			errors, token := libs.NewTokenWithTime(acs)
			if !errors {
				errs = false
				message = "Here the token!"
				data = token
			}
		} else {
			message = "Machine not authorized!"
		}
	} else {
		message = "One or more field empty!"
	}

	req := request{
		Error:   errs,
		Message: message,
		Token:   data,
	}
	return req
}
