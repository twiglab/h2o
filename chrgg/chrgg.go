package chrgg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

const CLIENT_ID = "chrgg"

type ChargeErr struct {
	Code    string
	Type    string
	Message string
}

func (e ChargeErr) Error() string {
	return fmt.Sprintf("Charge Error: code = %s, type = %s, message = %s", e.Code, e.Type, e.Message)
}

type ChargeData struct {
	MeterData
}

type ChangeServer struct {
	dbx    *DBx
	cdrLog *slog.Logger

	re RulerEngine
}

func (s *ChangeServer) pre(_ context.Context, md MeterData) (ChargeData, error) {
	return ChargeData{MeterData: md}, nil
}

func (s ChangeServer) check(ctx context.Context, last *ent.CDR, cd ChargeData) error {
	if cd.DataCode == last.DataCode {
		return &ChargeErr{Message: "same datacode"}
	}

	if !cd.DataTime.After(last.DataTime) {
		return &ChargeErr{Message: "time"}
	}
	return nil
}

func (s *ChangeServer) Verify(ctx context.Context, last *ent.CDR, cd ChargeData) bool {
	if isInTimeRange(cd.DataTime) {
		return true
	}

	if cd.Data.DataValue-last.Value < 100 { // 小于一个读数
		return false
	}
	return true
}

func (s *ChangeServer) DoChange(ctx context.Context, bd MeterData) (CDR, error) {
	// setp 1 prepare
	cd, err := s.pre(ctx, bd)
	if err != nil {
		return Nil, err
	}

	// step 2 load
	last, notfound, err := s.dbx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return Nil, err
	}

	// step 2.1 save first
	if notfound {
		nc := FirstCDR(cd)
		_, err = s.dbx.SaveCurrent(ctx, nc)
		return nc, err
	}

	// step 3 verify and check
	if !s.Verify(ctx, last, cd) {
		return Nil, nil
	}

	if err = s.check(ctx, last, cd); err != nil {
		return Nil, err
	}

	// setp 4 calc
	ru, err := s.re.GetResult(ctx, cd)
	if err != nil {
		return Nil, err
	}

	nc := CalcCDR(last, cd, ru)

	// step 5 write cdr

	// step 6 save
	_, err = s.dbx.SaveCurrent(ctx, nc)
	return nc, err
}

func isInTimeRange(t time.Time) bool {
	h, m, _ := t.Clock()
	if h*60+m > 1365 { // > 22:45
		return true
	}
	return false
}
