package common

import ()

func StringIn(str string, strList []string) (isIn bool) {
	for _, s := range strList {
		if s == str {
			return true
		}
	}
	return false
}
