package filter

import (
	"errors"
	"fmt"
	"io"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/ibbd-dev/go-csv"
)

func Filter(r io.Reader, w io.Writer, fieldnames []string, cmp func(str string) bool) (err error) {
	// 获取字段名
	reader, header, err := common.CsvMapReader(r)
	if err != nil {
		return err
	}

	if len(fieldnames) > 0 {
		for _, f := range fieldnames {
			if !common.StringIn(f, header) {
				return errors.New(f + " 该字段不存在")
			}
		}
	}

	writer := goCsv.NewMapWriter(w, header)
	if err = writer.WriteHeader(); err != nil {
		return err
	}

	var i int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		isFilter := false
		if len(fieldnames) > 0 {
			for _, f := range fieldnames {
				if cmp(record[f]) {
					// 如果满足条件则直接过滤
					//print(record[f], "====")
					isFilter = true
					break
				}
			}
		} else {
			for _, v := range record {
				if cmp(v) {
					isFilter = true
					break
				}
			}
		}

		if isFilter {
			continue
		}

		i++
		if common.Debug && i%1000 == 0 {
			fmt.Printf("Parse %d...\n", i)
		}

		if err = writer.WriteRow(record); err != nil {
			return err
		}

		if i%10 == 0 {
			writer.Flush()
		}
	}

	writer.Flush()
	return nil
}
