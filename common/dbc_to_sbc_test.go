package common

import (
	"testing"
)

func TestDBC2SBC(t *testing.T) {
	s1 := "。.，（）-1！@234567890abc１２３４５６７８９ａｂｃＡＢＣ　中文"
	s2 := "..,()-1!@234567890abc123456789abcABC 中文"
	if s2 != DBC2SBC(s1) {
		t.Fatalf("全角转半角错误：\"%s\" => \"%s\"", DBC2SBC(s1), s2)
	}
}
