package common

import "time"

func GetYear() int {
	return time.Now().Year()
}
