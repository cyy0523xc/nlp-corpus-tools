package corpus

import (
	"encoding/csv"
	"io"

	"github.com/cyy0523xc/nlp-corpus-tools/common"
)

func Batch(r io.Reader, ws []io.Writer, rate int, per int) (err error) {
	// 获取字段名
	reader, header, err := common.CsvReader(r)
	if err != nil {
		return err
	}

	wLen := len(ws)
	var writers = make([]*csv.Writer, wLen)
	for i, w := range ws {
		writers[i] = csv.NewWriter(w)
		if err = writers[i].Write(header); err != nil {
			return err
		}
	}

	rand := common.NewRand()
	var i, luck, wi int
	var wCount int // 记录每个文件写入的数量
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if rate > 0 {
			if i%rate == 0 {
				// 产生中奖号码
				luck = rand.Intn(rate)
			}
			if i%rate != luck {
				i++
				continue
			}
		}

		// 写入测试集
		current_wi := wi % wLen
		wi++
		if err = writers[current_wi].Write(record); err != nil {
			return err
		}
		writers[current_wi].Flush()

		if current_wi == wLen-1 {
			wCount++
			if wCount == per {
				return nil
			}
		}

		i++
	}
	return nil
}
