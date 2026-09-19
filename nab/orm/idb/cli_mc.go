package idb

import "modernc.org/sqlite/vtab"

var client_table_columns = []string{
	"id",
	"code",
	"typ",
	"url",
	"speed",
	"data_bits",
	"parity",
	"stop_bits",
	"timeout",
	"endian",
	"word_order",
	"memo",
}

type ClientRec struct {
	ID   int64  `csv:"id" db:"id"`
	Code string `csv:"code" db:"code"`
	Type string `csv:"type" db:"typ"`

	URL string `csv:"url" db:"url"`

	Speed    int64 `csv:"speed" db:"speed"`
	DataBits int64 `csv:"data_bits" db:"data_bits"`
	Parity   int64 `csv:"parity" db:"parity"`
	StopBits int64 `csv:"stop_bits" db:"stop_bits"`

	Timeout int64 `csv:"timeout" db:"timeout"`

	Endian    int64 `csv:"endian" db:"endian"`
	WordOrder int64 `csv:"word_order" db:"word_order"`

	Memo string `csv:"memo" db:"memo"`
}

func (r ClientRec) Column(i int) (vtab.Value, error) {
	switch i {
	case 0:
		return r.ID, nil
	case 1:
		return r.Code, nil
	case 2:
		return r.Type, nil
	case 3:
		return r.URL, nil
	case 4:
		return r.Speed, nil
	case 5:
		return r.DataBits, nil
	case 6:
		return r.Parity, nil
	case 7:
		return r.StopBits, nil
	case 8:
		return r.Timeout, nil
	case 9:
		return r.Endian, nil
	case 10:
		return r.WordOrder, nil
	case 11:
		return r.Memo, nil
	}
	panic(i)
}

func (c ClientRec) RowID() int64 {
	return c.ID
}
