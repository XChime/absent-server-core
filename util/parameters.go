/*
 * Copyright (c) 2019.
 * Davin Alfarizky Putra Basudewa  & dvnlabs.ml
 * dbasudewa@gmail.com
 * This is source for Finance API,not for public!
 */

package util

import "bytes"

type Parameters struct {
	Key   string
	Value string
}

/*Return error = false if no error*/
func CheckParameter(param []Parameters) (bool, string) {
	var buffer bytes.Buffer
	missing := 0
	buffer.WriteString("Field: ")
	for index, element := range param {
		if element.Value == "" && index < len(param) {
			missing++
			buffer.WriteString(element.Key)
		}
		if index+1 == len(param) && missing != 0 {
			buffer.WriteString(" Is missing!")
		}
		if element.Value == "" && index < len(param)-2 {
			buffer.WriteString(",")
		}
	}
	if missing != 0 {
		return true, buffer.String()
	}
	return false, ""
}
