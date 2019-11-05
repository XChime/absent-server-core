package main

import (
	"absensi-server/core/action/absent"
	"absensi-server/core/action/overtime"
	"absensi-server/core/action/schedule"
	"absensi-server/core/master/machine"
	"absensi-server/core/master/user"
	"absensi-server/libs/database"
	"absensi-server/util/data"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	inits := database.Connect()
	if inits != nil {
		initRoute()
	}
}

/*initRoute function untuk menginisialisasi route web*/
func initRoute() {
	//Getting some system key value of PORT
	var port = os.Getenv("PORT")
	if port == "" {
		println("Using default port 8080")
		port = "8080"
	}
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/login/administrator", administratorHandler).Methods("POST")
	//employee
	router.HandleFunc("/login/employee", employeeHandler).Methods("POST")
	router.HandleFunc("/create/employee", employeeCreateHandler).Methods("POST")
	router.HandleFunc("/reset/employee", employeeResetHandler).Methods("POST")
	router.HandleFunc("/change/employee/password", employeeChangePassHandler).Methods("POST")
	router.HandleFunc("/employee/div/{division}", employeeByDivision).Methods("GET")

	//machine
	router.HandleFunc("/login/machine", machineLoginHandler).Methods("POST")
	router.HandleFunc("/machine/request", machineRequestCodeHandler).Methods("POST")

	//Schedule
	//Create schedule for division
	router.HandleFunc("/employee/schedule/create", scheduleCreateHandler).Methods("POST")
	//Attach schedule to specific employee
	router.HandleFunc("/employee/schedule/attach", scheduleAttachHandler).Methods("POST")
	//Update schedule
	router.HandleFunc("/employee/schedule/update", scheduleUpdateHandler).Methods("POST")
	//Show specific employee schedule
	router.HandleFunc("/employee/schedule", scheduleEmployeeHandler).Methods("POST")
	//Show specific schedule by id
	router.HandleFunc("/employee/schedule/{id}", scheduleEmployeeById).Methods("GET")
	//Show schedule by division which is valid / not (GRANTED / WAITING)
	router.HandleFunc("/employee/schedule/{div}/{grant}", scheduleEmployeeByDiv).Methods("GET")

	//Overtime
	//Create overtime
	router.HandleFunc("/overtime/create", createOvertimeHandler).Methods("POST")
	//Update overtime
	router.HandleFunc("/overtime/update", updateOvertimeHandler).Methods("POST")
	//List overtime
	router.HandleFunc("/overtime/list/{div}/{grant}", listOvertimeDivision).Methods("GET")

	//OutDetails
	//Create OutDetail //TODO
	router.HandleFunc("/detail/out/create", nil).Methods("POST")
	//Update OutDetail //TODO
	router.HandleFunc("/detail/out/update", nil).Methods("POST")
	//Show Outdetail by division (GRANTED / WAITING)
	router.HandleFunc("/detail/out/show/{div}/{grant}", nil).Methods("POST")

	//Absent
	//Absent Request
	router.HandleFunc("/absent/request", requestAbsentHandler).Methods("POST")
	//Read Absent specific day
	router.HandleFunc("/absent/show/{day}", absentDayHandler).Methods("GET")
	//Read absent employee
	router.HandleFunc("/absent/show/employee/{nik}", absentEmployeeHandler).Methods("GET")

	// Handle all preflight request
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})
	router.StrictSlash(true)
	color.Warn.Println("Connected to port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	type homes struct {
		Message     string
		Description string
	}
	home := homes{
		Message:     "2019 (c) Kelompok 2 - All rights reserved",
		Description: "Absensi Server written by Davin Alfarizky Putra Basudewa",
	}
	var homeJson = string(data.MustMarshal(home))
	_, _ = fmt.Fprint(w, homeJson)
}
func employeeByDivision(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	params := mux.Vars(r)
	div := params["division"]

	emp := user.EmployeeByDivison(div)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}
func administratorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	nik := r.FormValue("nik")
	password := r.FormValue("password")

	emp := user.EmpLoginAdminHub(nik, password)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	nik := r.FormValue("nik")
	password := r.FormValue("password")
	deviceid := r.FormValue("device")

	emp := user.EmpHubLogin(nik, password, deviceid)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}

func employeeChangePassHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	nik := r.FormValue("nik")
	password := r.FormValue("password")

	emp := user.ChangePasswordHub(nik, password)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}

func employeeResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	nik := r.FormValue("nik")
	token := r.FormValue("token")

	emp := user.ResetAccountHub(nik, token)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}

func employeeCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	nameCreate := r.FormValue("nameCreate")
	divisi := r.FormValue("divisi")
	token := r.FormValue("token")

	emp := user.CreateHub(nameCreate, divisi, token)
	var homeJson = string(data.MustMarshal(emp))
	_, _ = fmt.Fprint(w, homeJson)
}

//Machine section

func machineLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	sharecode := r.FormValue("sharecode")
	mch := machine.LoginMachineHub(sharecode)
	var homeJson = string(data.MustMarshal(mch))
	_, _ = fmt.Fprint(w, homeJson)
}

func machineRequestCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	id := r.FormValue("id")
	secret := r.FormValue("secret")
	mch := machine.RequestMachineAccessHub(id, secret)
	var homeJson = string(data.MustMarshal(mch))
	_, _ = fmt.Fprint(w, homeJson)
}

//Schedule section
func scheduleCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	divisi := r.FormValue("divisi")
	mesg := r.FormValue("message")
	datas := r.FormValue("data")
	sch := schedule.CreateScheduleHub(token, divisi, mesg, datas)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func scheduleAttachHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	idsch := r.FormValue("schedule")
	nik := r.FormValue("employee")
	sch := schedule.ScheduleAttachHub(idsch, token, nik)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func scheduleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	idsch := r.FormValue("schedule")
	message := r.FormValue("message")
	datas := r.FormValue("data")
	grant := r.FormValue("grant")
	sch := schedule.ScheduleUpdateHub(idsch, token, message, datas, grant)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func scheduleEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	sch := schedule.ReadEmployeeScheduleHub(token)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}
func scheduleEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	params := mux.Vars(r)
	id := params["id"]
	sch := schedule.ReadSchedulebyIdHub(id)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}
func scheduleEmployeeByDiv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)

	params := mux.Vars(r)
	div := params["div"]
	grant := params["grant"]
	divI, _ := strconv.Atoi(div)
	sch := schedule.ScheduleByDivisionHub(divI, grant)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

/*OVERTIME
FUNCTION*/
func createOvertimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	offset := r.FormValue("offset")
	message := r.FormValue("message")
	divisi := r.FormValue("division")
	offsetI, _ := strconv.Atoi(offset)
	divisiI, _ := strconv.Atoi(divisi)
	sch := overtime.CreateOvertimeHub(token, offsetI, message, divisiI)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func updateOvertimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	token := r.FormValue("token")
	id := r.FormValue("overtime")
	offset := r.FormValue("offset")
	message := r.FormValue("message")
	divisi := r.FormValue("division")
	status := r.FormValue("status")
	offsetI, _ := strconv.Atoi(offset)
	divisiI, _ := strconv.Atoi(divisi)
	sch := overtime.UpdateOvertimeHub(token, id, offsetI, message, divisiI, status)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func listOvertimeDivision(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)

	params := mux.Vars(r)
	div := params["div"]
	grant := params["grant"]
	divI, _ := strconv.Atoi(div)
	sch := overtime.ReadOvertimeListHub(divI, grant)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)

}

/*Absent Section
2019 11 02*/
func requestAbsentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)
	if err := r.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	tknclient := r.FormValue("client")
	tknmachine := r.FormValue("machine")
	deviceid := r.FormValue("deviceid")

	sch := absent.RequestAbsentHub(tknclient, tknmachine, deviceid)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func absentEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)

	params := mux.Vars(r)
	nik := params["nik"]
	sch := absent.ReadAbsentByEmployeeHub(nik)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func absentDayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	LogConsoleHttpReq(r)

	params := mux.Vars(r)
	day := params["day"]
	sch := absent.ReadAbsentByDaysHub(day)
	var homeJson = string(data.MustMarshal(sch))
	_, _ = fmt.Fprint(w, homeJson)
}

func LogConsoleHttpReq(r *http.Request) {
	color.Cyan.Println(r.Method + " : " + r.Proto + " [" + r.Host + r.URL.String() +
		"] Requested by: " + r.RemoteAddr + " At:->" + time.Now().Format("Mon, 2 Jan 2006 15:04:05 MST"))
}
