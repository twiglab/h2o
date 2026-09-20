package idb

import (
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
)

func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.TrimLeadingSpace = true
		r.Comment = '#'
		return r
	})
}
