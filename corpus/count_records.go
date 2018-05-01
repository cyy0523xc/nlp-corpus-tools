package corpus

import (
	"fmt"
	"io"
	"strings"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/ibbd-dev/go-csv"
)

func Stat(r io.Reader) (n int, err error) {
	reader := goCsv.NewReader(r)
	if header, err := reader.GetFieldnames(); err != nil {
		return n, err
	} else {
		fmt.Printf("字段数：%d\n字段：%+v\n", len(header), strings.Join(header, ", "))
	}

	// 判断编码
	var isUTF8 = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, nil
		}

		if isUTF8 {
			str := strings.Join(record, ", ")
			if !common.IsUTF8([]byte(str)) {
				isUTF8 = false
			}
		}

		n++
	}

	if isUTF8 {
		fmt.Println("该文件是utf8编码")
	} else {
		fmt.Println("该文件不是utf8编码，注意判断")
	}
	return n, nil
}
