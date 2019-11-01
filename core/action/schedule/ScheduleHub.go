package schedule

import (
	"absensi-server/core/master/user/model"
	"absensi-server/libs"
	deta "absensi-server/util/data"
	"encoding/json"
)

func CreateScheduleHub(token string, divisi string, mesg string, datas string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error   bool
		Message string
		Data    interface{}
	}
	if token != "" && divisi != "" && mesg != "" && datas != "" {
		isOK, tok, mssg := libs.VerifyToken(token)
		if isOK {
			jsonString := deta.MustMarshal(tok)
			var admin model.AdminData
			_ = json.Unmarshal(jsonString, &admin)
			isOKk, msg, dat := CreateSchedule(admin.NIK, divisi, mesg, datas)
			if isOKk {
				data = dat
				errs = false
			}
			message = msg
		} else {
			message = mssg
		}
	} else {
		message = "One or more field is empty!"
	}

	sch := schedule{
		Error:   errs,
		Message: message,
		Data:    data,
	}
	return sch
}

func ReadSchedulebyIdHub(id string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error   bool
		Message string
		Data    interface{}
	}

	if id != "" {
		dats, msg := ScheduleById(id)
		if dats != nil {
			errs = false
			data = dats
			message = msg
		} else {
			message = msg
		}
	}

	sch := schedule{
		Error:   errs,
		Message: message,
		Data:    data,
	}
	return sch
}

func ScheduleByDivisionHub(division int, grant string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error   bool
		Message string
		Data    interface{}
	}

	if division != 0 && grant != "" {
		dats, msg := ScheduleDivison(division, grant)
		if dats != nil {
			errs = false
			data = dats
			message = msg
		} else {
			message = "No data"
		}
	}

	sch := schedule{
		Error:   errs,
		Message: message,
		Data:    data,
	}
	return sch

}
func ScheduleAttachHub(idsch string, token string, nik string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error    bool
		Message  string
		Schedule interface{}
	}

	if idsch != "" && token != "" && nik != "" {
		isOK, tok, mssg := libs.VerifyToken(token)
		if isOK {
			jsonString := deta.MustMarshal(tok)
			var admin model.AdminData
			_ = json.Unmarshal(jsonString, &admin)
			if deta.IsNIKAdmin(admin.NIK) && admin.IsAdmin == true {
				dats, msg := ScheduleAttach(idsch, nik)
				if dats != nil {
					data = dats
					errs = false
				}
				message = msg
			} else {
				message = "Unsufficient Access!"
			}
		} else {
			message = mssg
		}
	} else {
		message = "One or more field not found!"
	}

	sch := schedule{
		Error:    errs,
		Message:  message,
		Schedule: data,
	}
	return sch
}

func ScheduleUpdateHub(idsch string, token string, msgs string, datas string, grant string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error    bool
		Message  string
		Schedule interface{}
	}
	if idsch != "" && token != "" && msgs != "" && datas != "" {
		isOK, tok, mssg := libs.VerifyToken(token)
		if isOK {
			jsonString := deta.MustMarshal(tok)
			var admin model.AdminData
			_ = json.Unmarshal(jsonString, &admin)
			if deta.IsNIKAdmin(admin.NIK) && admin.IsAdmin == true {
				dats, msg := ScheduleUpdate(idsch, admin.NIK, msgs, datas, grant)
				if dats != nil {
					data = dats
					errs = false
				}
				message = msg
			} else {
				message = "Unsufficient Access!"
			}
		} else {
			message = mssg
		}

	} else {
		message = "One or more field not entered!"
	}
	sch := schedule{
		Error:    errs,
		Message:  message,
		Schedule: data,
	}
	return sch
}

func ReadEmployeeScheduleHub(token string) interface{} {
	errs := true
	message := ""
	var data interface{}
	type schedule struct {
		Error   bool
		Message string
		Data    interface{}
	}

	if token != "" {
		isOk, dat, msgs := libs.VerifyToken(token)
		if isOk {
			jsonString := deta.MustMarshal(dat)
			var loginData model.LoginData
			_ = json.Unmarshal(jsonString, &loginData)
			dat, msg := ReadEmployeeSchedule(loginData.NIK)
			if dat != nil {
				data = dat
				errs = false
				message = msg
			} else {
				message = msg
			}
		} else {
			message = msgs
		}

	} else {
		message = "Token can't find!"
	}

	schedules := schedule{
		Error:   errs,
		Message: message,
		Data:    data,
	}
	return schedules

}
