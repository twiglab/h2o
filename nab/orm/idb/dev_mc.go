package idb

import (
	"fmt"

	"modernc.org/sqlite/vtab"
)

var DEVICE_TABLE_COLUMNS = []string{
	"id",
	"code",
	"typ",
	"sn",
	"clazz",
	"cli",
	"unit_id",
	"addr",
	"endian",
	"word_order",
	"memo",
}

type DevRec struct {
	ID   int64  `csv:"id" db:"id"`
	Code string `csv:"code" db:"code"`
	Typ  string `csv:"typ" db:"typ"`
	SN   string `csv:"sn" db:"sn"`

	Clazz string `csv:"clazz" db:"clazz"`

	Cli       string `csv:"cli" db:"cli"`
	UnitID    int64  `csv:"unit_id" db:"unit_id"`
	Addr      int64  `csv:"addr" db:"addr"`
	Endian    int64  `csv:"endian" db:"endian"`
	WordOrder int64  `csv:"word_order" db:"word_order"`

	Memo string `csv:"memo" db:"memo"`
}

func (r DevRec) Column(i int) (vtab.Value, error) {
	switch i {
	case 0:
		return r.ID, nil
	case 1:
		return r.Code, nil
	case 2:
		return r.Typ, nil
	case 3:
		return r.SN, nil
	case 4:
		return r.Clazz, nil
	case 5:
		return r.Cli, nil
	case 6:
		return r.UnitID, nil
	case 7:
		return r.Addr, nil
	case 8:
		return r.Endian, nil
	case 9:
		return r.WordOrder, nil
	case 10:
		return r.Memo, nil
	}
	panic(fmt.Errorf("no field %d", i))
}

func (r DevRec) RowID() int64 {
	return r.ID
}
