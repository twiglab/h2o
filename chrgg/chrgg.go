package chrgg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

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

func (s *ChangeServer) calc(ctx context.Context, cd ChargeData) (CDR, error) {
	last, notfound, err := s.dbx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return Nil, err
	}

	if notfound {
		return FirstCDR(cd), nil
	}

	if err := s.check(ctx, last, cd); err != nil {
		return Nil, err
	}

	if !s.Verify(ctx, last, cd) {
		return Nil, nil
	}

	ru, err := s.re.GetResult(ctx, cd)
	if err != nil {
		return Nil, err
	}

	return CalcCDR(last, cd, ru), nil
}

func (s *ChangeServer) DoChange(ctx context.Context, bd MeterData) (nc CDR, err error) {
	var cd ChargeData
	if cd, err = s.pre(ctx, bd); err != nil {
		return
	}

	nc, err = s.calc(ctx, cd)
	if err != nil {
		return
	}
	_, err = s.dbx.SaveCurrent(ctx, nc)
	return
}

func isInTimeRange(t time.Time) bool {
	start := time.Date(t.Year(), t.Month(), t.Day(), 22, 15, 0, 0, t.Location())
	end := start.Add(time.Hour)
	return !t.Before(start) && !t.After(end)
}
