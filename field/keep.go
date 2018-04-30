package field

import (
	"errors"
	"fmt"
	"io"

	//"github.com/cyy0523xc/nlp-corpus-tools/common"
	"github.com/ibbd-dev/go-csv"
)

func Keep(r io.Reader, w io.Writer, fieldnames []string) (err error) {
	reader := goCsv.NewMapReader(r)
	writer := goCsv.NewMapWriterSimple(w)
	writer.SetHeader(fieldnames)
	writer.WriteHeader()

	var i = 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if i == 0 {
			for _, f := range fieldnames {
				if _, ok := record[f]; !ok {
					return errors.New(fmt.Sprintf("%s 字段名不存在, Line: %d", f, i))
				}
			}
		}

		i++
		writer.WriteRow(record)
		if i%10 == 0 {
			writer.Flush()
		}
	}

	writer.Flush()
	return nil
}
