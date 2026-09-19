package idb

import (
	"modernc.org/sqlite/vtab"
)

type VTable struct{}

func (VTable) BestIndex(info *vtab.IndexInfo) error { return nil }
func (VTable) Open() (vtab.Cursor, error)           { return nil, nil }
func (VTable) Disconnect() error                    { return nil }
func (VTable) Destroy() error                       { return nil }

type VitrualRecord interface {
	RowID() int64
	Column(col int) (vtab.Value, error)
}

type VirtualTable[R VitrualRecord] struct {
	Rows []R
}

func (VirtualTable[T]) BestIndex(info *vtab.IndexInfo) error { return nil }
func (t VirtualTable[T]) Open() (vtab.Cursor, error) {
	return &VirtualCursor[T]{
		Pos:   0,
		Rows:  t.Rows,
		Table: t,
	}, nil
}
func (VirtualTable[T]) Disconnect() error { return nil }
func (VirtualTable[T]) Destroy() error    { return nil }

/*
func (VirtualTable) Columns() []string { return nil }

func (VirtualTable) Insert(cols []vtab.Value, rowid *int64) error                    { return nil }
func (VirtualTable) Update(oldRowid int64, cols []vtab.Value, newRowid *int64) error { return nil }
func (VirtualTable) Delete(oldRowid int64) error                                     { return nil }
*/

type VCursor struct{}

func (VCursor) Filter(_ int, idxStr string, _ []vtab.Value) error { return nil }
func (VCursor) Next() error                                       { return nil }
func (VCursor) Column(_ int) (vtab.Value, error)                  { return nil, nil }
func (VCursor) Eof() bool                                         { return false }
func (VCursor) Rowid() (int64, error)                             { return 0, nil }
func (VCursor) Close() error                                      { return nil }

type VirtualCursor[R VitrualRecord] struct {
	Pos   int
	Rows  []R
	Table vtab.Table
}

func (VirtualCursor[R]) Filter(idxNum int, idxStr string, vals []vtab.Value) error { return nil }

func (c *VirtualCursor[R]) Next() error {
	if c.Pos < len(c.Rows) {
		c.Pos++
	}
	return nil
}

func (c VirtualCursor[R]) Column(col int) (vtab.Value, error) {
	if c.Pos >= len(c.Rows) {
		return nil, nil
	}
	return c.Rows[c.Pos].Column(col)
}
func (c VirtualCursor[R]) Eof() bool {
	return c.Pos >= len(c.Rows)
}
func (c VirtualCursor[R]) Rowid() (int64, error) {
	return c.Rows[c.Pos].RowID(), nil
}
func (VirtualCursor[R]) Close() error { return nil }
