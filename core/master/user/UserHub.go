package user

import (
	"absensi-server/core/master/user/model"
	"absensi-server/libs"
	"absensi-server/util/common"
	deta "absensi-server/util/data"
	"encoding/json"
)

//Function untuk user/karyawan baru
func CreateHub(nameCreate string, divisi string, token string) interface{} {
	errors := true
	var data interface{}
	var message string
	type EmpCreate struct {
		Error   bool
		Message string
		Data    interface{}
	}
	if token != "" && nameCreate != "" && divisi != "" {
		isOk, dat, msgs := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.LoginData
			_ = json.Unmarshal(jsonString, &loginData)
			if deta.IsNIKAdmin(loginData.NIK) {
				datas, msg := Create(nameCreate, divisi)
				if datas != nil {
					data = datas
					message = msg
					errors = false
				} else {
					message = msg
					errors = true
				}
			} else {
				errors = true
				message = "Required Personalia / IT Department access"
			}
		} else {
			errors = true
			message = msgs
		}
	} else {
		errors = true
		message = "One or more field empty!"
	}

	created := EmpCreate{
		Error:   errors,
		Message: message,
		Data:    data,
	}
	return created
}

func EmpHubLogin(nik string, password string, deviceId string) interface{} {
	errors := true
	var data interface{}
	var message string
	type EmpAuth struct {
		Error   bool
		Message string
		Data    interface{}
	}
	if !common.VarStringChecker(nik) {
		message = "NIK empty"
	} else if !common.VarStringChecker(password) {
		message = "Password empty"
	} else if common.VarStringChecker(nik) && common.VarStringChecker(password) && deviceId != "" {
		errors = false
		LoginError, datas, msg := Login(nik, password, deviceId)
		if !LoginError {
			errStatus, token := libs.NewToken(datas)
			if errStatus {
				errors = true
				msg = "JWT ERROR"
			} else {
				data = token
				message = msg
			}
		} else {
			errors = true
			message = msg
		}

	} else if !common.VarStringChecker(nik) && !common.VarStringChecker(password) {
		errors = true
		message = "NIK & Password empty"
	}

	authorized := EmpAuth{
		Error:   errors,
		Message: message,
		Data:    data,
	}
	return authorized

}

func ResetAccountHub(nik string, token string) interface{} {
	type EmpReset struct {
		Error   bool
		Message string
		Data    interface{}
	}
	var data interface{}
	var message = ""
	erroor := true
	if nik != "" && token != "" {
		isOk, dat, msgs := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.LoginData
			_ = json.Unmarshal(jsonString, &loginData)
			if deta.IsNIKAdmin(loginData.NIK) {
				datas, msg := ResetAccount(nik)
				if datas != nil {
					erroor = false
					data = datas
					message = msg
				} else {
					erroor = true
					message = msg
				}
			} else {
				erroor = true
				message = "Required Personalia / IT Department access"
			}
		} else {
			erroor = true
			message = msgs
		}
	} else {
		erroor = true
		message = "One or more field empty!"
	}
	reseter := EmpReset{
		Error:   erroor,
		Message: message,
		Data:    data,
	}
	return reseter
}
func ChangePasswordHub(nik string, password string) interface{} {
	erroor := true
	message := ""
	type changed struct {
		Error   bool
		Message string
	}
	if nik != "" && password != "" {
		ok, msg := ChangePassword(nik, password)
		if ok {
			erroor = false
			message = msg
		} else {
			erroor = true
			message = msg
		}
	} else {
		message = "One or more field required!"
	}

	changes := changed{
		Error:   erroor,
		Message: message,
	}
	return changes
}
