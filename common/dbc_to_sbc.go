package common

import ()

func DBC2SBC(s string) string {
	retstr := ""
	for _, i := range s {
		inside_code := i
		if i == 12290 {
			inside_code = 46
		} else if i == 12288 { // 空格
			inside_code = 32
		} else {
			inside_code -= 65248
		}

		//print(string(i), "  ", i, " ", string(inside_code), "\n")
		if inside_code < 32 || inside_code > 126 {
			retstr += string(i)
		} else {
			retstr += string(inside_code)
		}
	}
	return retstr
}
