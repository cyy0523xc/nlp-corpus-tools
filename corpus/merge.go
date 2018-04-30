package corpus

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
)

func Merge(rs []io.Reader, w io.Writer) (err error) {
	writer := csv.NewWriter(w)

	var reader *csv.Reader
	for i := range rs {
		reader = csv.NewReader(rs[i])
		if i == 0 {
			if header, err := reader.Read(); err != nil {
				return err
			} else {
				// 只需要写入一次header
				if err = writer.Write(header); err != nil {
					return err
				}
			}
		} else {
			if _, err = reader.Read(); err != nil {
				return err
			}
		}

		if common.Debug {
			fmt.Printf("开始第%d个文件读取...\n", i)
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if err = writer.Write(record); err != nil {
				return err
			}
			writer.Flush()
		}
	}

	return nil
}
