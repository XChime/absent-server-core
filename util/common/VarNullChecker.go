package common

//VarStringChecker return true if string not empty
func VarStringChecker(checks string) bool {
	if checks != "" {
		return true
	}
	return false
}
