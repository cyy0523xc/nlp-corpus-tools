package corpus

import (
	"io"

	"github.com/ibbd-dev/go-csv"
)

func CountRecords(r io.Reader) (n int, err error) {
	return goCsv.CountLines(r)
}
