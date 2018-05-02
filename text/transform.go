package text

import (
	"errors"
	"fmt"
	"io"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/ibbd-dev/go-csv"
)

func Transform(r io.Reader, w io.Writer, fieldnames []string, function func(str string) (string, error)) (err error) {
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

		i++
		if common.Debug && i%1000 == 0 {
			fmt.Printf("Parse %d...\n", i)
		}

		if len(fieldnames) > 0 {
			for _, f := range fieldnames {
				if record[f], err = function(record[f]); err != nil {
					return err
				}
			}
		} else {
			for f, v := range record {
				if record[f], err = function(v); err != nil {
					return err
				}
			}
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
