package machine

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
