package common

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	GBK     = "gbk"
	GB18030 = "gb18030"
)

var (
	// 默认为gbk编码
	chineseGB = simplifiedchinese.GBK
)

// IsUTF8 判断是否为utf8编码
func IsUTF8(bts []byte) bool {
	return utf8.Valid(bts)
}

// SetGBCharset 可以改变默认的中文编码
func SetGBCharset(typ string) error {
	typ = strings.ToLower(typ)
	switch typ {
	case GBK:
		chineseGB = simplifiedchinese.GBK
	case GB18030:
		chineseGB = simplifiedchinese.GB18030
	default:
		return errors.New("not support charset: " + typ)
	}

	return nil
}

// GB_UTF8 中文编码转化为utf8编码
func GB_UTF8(s string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), chineseGB.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

// UTF8_GB utf8编码转化为中文编码
func UTF8_GB(s string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), chineseGB.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}
