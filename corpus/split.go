package corpus

import (
	"encoding/csv"
	"io"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
)

func Split(r io.Reader, train_w, test_w io.Writer, rate int) (err error) {
	// 获取字段名
	reader, header, err := common.CsvReader(r)
	if err != nil {
		return err
	}

	train_writer := csv.NewWriter(train_w)
	if err = train_writer.Write(header); err != nil {
		return err
	}
	test_writer := csv.NewWriter(test_w)
	if err = test_writer.Write(header); err != nil {
		return err
	}

	rand := common.NewRand()
	var i, luck int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if i%rate == 0 {
			// 产生中奖号码
			luck = rand.Intn(rate)
			train_writer.Flush()
			test_writer.Flush()
		}
		if i%rate == luck {
			// 写入测试集
			if err = test_writer.Write(record); err != nil {
				return err
			}
		} else {
			// 写入训练集
			if err = train_writer.Write(record); err != nil {
				return err
			}
		}

		i++
	}

	train_writer.Flush()
	test_writer.Flush()
	return
}
