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
			var loginData model.AdminData
			_ = json.Unmarshal(jsonString, &loginData)
			isADMINDB := deta.IsNIKAdmin(loginData.NIK)
			if isADMINDB && loginData.IsAdmin == true {
				datas, msg := Create(nameCreate, divisi)
				if datas != nil {
					data = datas
					message = msg
					errors = false
				} else {
					message = msg
					errors = true
				}
			} else if isADMINDB {
				errors = true
				message = "Required administrator token"
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
func EmployeeByDivison(div string) interface{} {
	type employee struct {
		Error    bool
		Message  string
		Employee interface{}
	}
	message := ""
	errs := true
	var data interface{}

	if div != "" {
		dats, msg := GetEmployeeByDivision(div)
		if dats != "" {
			errs = false
			data = dats
		}
		message = msg
	}

	emp := employee{
		Error:    errs,
		Message:  message,
		Employee: data,
	}
	return emp
}
func EmpLoginAdminHub(nik string, password string) interface{} {
	type admin struct {
		Error   bool
		Message string
		Admin   interface{}
	}
	message := ""
	errs := true
	var data interface{}

	if nik != "" && password != "" {
		isErr, dats, msg := LoginAdmin(nik, password)
		errs = isErr
		if !isErr {
			errStatus, token := libs.NewToken(dats)
			if errStatus {
				message = "JWT Error"
			} else {
				data = token
				message = msg
			}
		}

	} else {
		message = "One or more field not found!"
	}

	admins := admin{
		Error:   errs,
		Message: message,
		Admin:   data,
	}
	return admins
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
			var loginData model.AdminData
			_ = json.Unmarshal(jsonString, &loginData)
			isADMINDB := deta.IsNIKAdmin(loginData.NIK)
			if isADMINDB && loginData.IsAdmin == true {
				datas, msg := ResetAccount(nik)
				if datas != nil {
					erroor = false
					data = datas
					message = msg
				} else {
					erroor = true
					message = msg
				}
			} else if isADMINDB {
				erroor = true
				message = "Required administrator token"
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

func ShowEmployeeProfileHub(token string) interface{} {
	type profile struct {
		Error   bool
		Message string
		Profile interface{}
	}
	errs := true
	msg := ""
	var data interface{}
	isOk, dat, msgs := libs.VerifyToken(token)
	if isOk {
		jsonString := deta.MustMarshal(dat)
		var loginData model.AdminData
		_ = json.Unmarshal(jsonString, &loginData)
		nik := loginData.NIK
		if nik != "" {
			dats, messg := ShowEmployeeProfile(nik)
			if dats != nil {
				data = dats
				errs = false
			}
			msg = messg
		}
	} else {
		msg = msgs
	}

	prof := profile{
		Error:   errs,
		Message: msg,
		Profile: data,
	}
	return prof
}
