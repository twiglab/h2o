package idb

import (
	"modernc.org/sqlite/vtab"
)

type emptyVirtualTable struct{}

func (emptyVirtualTable) Columns() []string                                               { return nil }
func (emptyVirtualTable) BestIndex(info *vtab.IndexInfo) error                            { return nil }
func (emptyVirtualTable) Open() (vtab.Cursor, error)                                      { return nil, nil }
func (emptyVirtualTable) Disconnect() error                                               { return nil }
func (emptyVirtualTable) Destroy() error                                                  { return nil }
func (emptyVirtualTable) Insert(cols []vtab.Value, rowid *int64) error                    { return nil }
func (emptyVirtualTable) Update(oldRowid int64, cols []vtab.Value, newRowid *int64) error { return nil }
func (emptyVirtualTable) Delete(oldRowid int64) error                                     { return nil }

type emptyVirtualCursor struct{}

func (emptyVirtualCursor) Filter(idxNum int, idxStr string, vals []vtab.Value) error { return nil }
func (emptyVirtualCursor) Next() error                                               { return nil }
func (emptyVirtualCursor) Column(col int) (vtab.Value, error)                        { return nil, nil }
func (emptyVirtualCursor) Eof() bool                                                 { return false }
func (emptyVirtualCursor) Rowid() (int64, error)                                     { return 0, nil }
func (emptyVirtualCursor) Close() error                                              { return nil }
