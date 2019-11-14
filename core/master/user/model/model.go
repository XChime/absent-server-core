package model

type CreatedEmployee struct {
	NIK             string
	DefaultPassword string
	Message         string
}

//For admin application
type AdminData struct {
	IsAdmin    bool
	NIK        string
	Nama       string
	DivisiName string
	Divisi     int
}

//For standard account
type LoginData struct {
	NIK        string
	Nama       string
	DivisiName string
	Divisi     int
	Jadwal     string
	DeviceHash string
}
type EmployeList struct {
	List []EmployeeListDetail
}
type EmployeeListDetail struct {
	NIK        string
	Nama       string
	DivisiName string
	Divisi     int
}

type ProfileEmployee struct {
	NIK        string
	Nama       string
	Divisi     int
	NamaDivisi string
	Jadwal     string
}
