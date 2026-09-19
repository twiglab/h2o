package idb

import (
	"errors"
	"fmt"
	"strings"

	"modernc.org/sqlite/vtab"
)

type IDBModule struct {
}

func (m *IDBModule) Create(ctx vtab.Context, args []string) (vtab.Table, error) {
	file := parseCSVArgs(args[3:])
	if file == "" {
		return nil, fmt.Errorf("csv: require filename=... arg")
	}

	switch args[2] {
	case "dev":
		return m.CreateDev(ctx, file)
	case "cli":
		return m.CreateCli(ctx, file)
	}

	return nil, errors.New("no table " + args[2])
}

func (m *IDBModule) CreateDev(ctx vtab.Context, file string) (vtab.Table, error) {
	t, err := loadTable[DevRec](file)
	if err != nil {
		return nil, err
	}
	if err := ctx.Declare(fmt.Sprintf("CREATE TABLE %s(%s)", "dev", strings.Join(DEVICE_TABLE_COLUMNS, ","))); err != nil {
		return nil, err
	}
	return t, nil
}

func (m *IDBModule) CreateCli(ctx vtab.Context, file string) (vtab.Table, error) {
	t, err := loadTable[CliRec](file)
	if err != nil {
		return nil, err
	}
	if err := ctx.Declare(fmt.Sprintf("CREATE TABLE %s(%s)", "cli", strings.Join(CLIENT_TABLE_COLUMNS, ","))); err != nil {
		return nil, err
	}
	return t, nil
}

func (m *IDBModule) Connect(ctx vtab.Context, args []string) (vtab.Table, error) {
	return m.Create(ctx, args)
}
