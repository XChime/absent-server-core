package main

import (
	"absensi-server/core/master/user"
	"absensi-server/libs/database"
	"absensi-server/util/data"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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
	//employee
	router.HandleFunc("/login/employee", employeeHandler).Methods("POST")
	router.HandleFunc("/create/employee", employeeCreateHandler).Methods("POST")
	router.HandleFunc("/reset/employee", employeeResetHandler).Methods("POST")
	router.HandleFunc("/change/employee/password", employeeChangePassHandler).Methods("POST")

	//machine
	//TODO adding machine

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

func LogConsoleHttpReq(r *http.Request) {
	color.Cyan.Println(r.Method + " : " + r.Proto + " [" + r.Host + r.URL.String() +
		"] Requested by: " + r.RemoteAddr + " At:->" + time.Now().String())
}
