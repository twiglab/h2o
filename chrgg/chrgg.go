package chrgg

import (
	"context"
	"fmt"
	"log/slog"

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
	if last.DataCode == cd.DataCode {
		return &ChargeErr{Message: "same datacode"}
	}

	if cd.DataTime.Compare(last.DataTime) < 1 {
		return &ChargeErr{Message: "time"}
	}
	return nil
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
	if err != nil { // err == nil
		return
	}
	_, err = s.dbx.SaveCurrent(ctx, nc)
	return
}
