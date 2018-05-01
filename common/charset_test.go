package common

import (
	"fmt"
	"testing"
)

func TestCharset(t *testing.T) {
	s := "GBK 与 UTF-8 编码转换测试"
	if !IsUTF8([]byte(s)) {
		t.Fatal("编码判断出错")
	}
	gbk, err := UTF8_GB(s)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(gbk)
	}
	if IsUTF8([]byte(gbk)) {
		t.Fatal("编码判断出错")
	}

	utf8, err := GB_UTF8(gbk)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(utf8)
	}
}
