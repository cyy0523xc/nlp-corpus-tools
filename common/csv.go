package common

import (
	"io"

	"github.com/ibbd-dev/go-csv"
)

func CsvReader(r io.Reader) (reader *goCsv.Reader, header []string, err error) {
	reader = goCsv.NewReader(r)
	if Skip > 0 {
		reader.SetSkip(Skip)
	}
	if Limit > 0 {
		reader.SetLimit(Limit)
	}
	header, err = reader.GetFieldnames()
	return
}

func CsvMapReader(r io.Reader) (reader *goCsv.MapReader, header []string, err error) {
	reader = goCsv.NewMapReader(r)
	if Skip > 0 {
		reader.SetSkip(Skip)
	}
	if Limit > 0 {
		reader.SetLimit(Limit)
	}
	header, err = reader.GetFieldnames()
	return
}
