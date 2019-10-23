package model

type CreatedEmployee struct {
	NIK             string
	DefaultPassword string
	Message         string
}

type Logined struct {
	Token string
}
type LoginData struct {
	NIK        string
	Divisi     int
	Jadwal     string
	DeviceHash string
}
