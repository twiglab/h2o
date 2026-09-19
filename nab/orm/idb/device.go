package idb

import (
	"modernc.org/sqlite/vtab"
)

var device_table_columns = []string{
	"id",
	"code",
	"typ",
	"sn",
	"clazz",
	"cli",
	"unit_id",
	"endian",
	"word_order",
}

type DeviceRec struct {
	ID   int64  `csv:"id" db:"id"`
	Code string `csv:"code" db:"code"`
	Type string `csv:"typ" db:"typ"`
	SN   string `csv:"sn" db:"sn"`

	Clazz string `csv:"clazz" db:"clazz"`

	Cli       string `csv:"cli" db:"cli"`
	UnitID    int64  `csv:"unit_id" db:"unit_id"`
	Endian    int64  `csv:"endian" db:"endian"`
	WordOrder int64  `csv:"word_order" db:"word_order"`
}

func (r DeviceRec) Column(i int) (vtab.Value, error) {
	switch i {
	case 0:
		return r.ID, nil
	case 1:
		return r.Code, nil
	case 2:
		return r.Type, nil
	case 3:
		return r.SN, nil
	case 4:
		return r.Clazz, nil
	case 5:
		return r.Cli, nil
	case 6:
		return r.UnitID, nil
	case 7:
		return r.Endian, nil
	case 8:
		return r.WordOrder, nil
	}
	panic("no filed")
}

func (r DeviceRec) RowID() int64 {
	return r.ID
}
