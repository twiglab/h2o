package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

type ChargeServer struct {
	DBx    *DBx
	CDRLog *slog.Logger
	Eng    ChargeEngine
}

func (s *ChargeServer) pre(_ context.Context, md MeterData) (ChargeData, error) {
	return ChargeData{MeterData: md}, nil
}

func (s ChargeServer) check(ctx context.Context, last *ent.CDR, cd ChargeData) error {
	if cd.DataCode == last.DataCode {
		return ErrDataCodeDup
	}

	if !cd.DataTime.After(last.DataTime) {
		return ErrTimeBefore
	}
	return nil
}

func (s *ChargeServer) Verify(ctx context.Context, last *ent.CDR, cd ChargeData) bool {
	if MinOfDay(hourMin(cd.DataTime)) >= 1365 {
		return true
	}

	if cd.Data.DataValue-last.Value < 100 { // 小于一个读数
		return false
	}
	return true
}

func (s *ChargeServer) doNewCDR(ctx context.Context, cd ChargeData) (CDR, error) {
	nc := FirstCDR(cd, ZeroRuler("new"))
	s.CDRLog.InfoContext(ctx, "cdr", slog.Any("cdr", nc))
	_, err := s.DBx.SaveCurrent(ctx, nc)
	return nc, err
}

func (s *ChargeServer) DoChange(ctx context.Context, bd MeterData) (CDR, error) {
	// setp 1 prepare
	cd, err := s.pre(ctx, bd)
	if err != nil {
		return Nil, err
	}

	// step 2 load
	last, notfound, err := s.DBx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return Nil, err
	}

	// step 2.1 save first
	if notfound {
		return s.doNewCDR(ctx, cd)
	}

	// step 3 verify and check
	if !s.Verify(ctx, last, cd) {
		return Nil, nil
	}

	if err = s.check(ctx, last, cd); err != nil {
		return Nil, err
	}

	// setp 4 calc
	ru, err := s.Eng.GetRuler(ctx, cd)
	if err != nil {
		return Nil, err
	}

	nc := CalcCDR(last, cd, ru)

	// step 5 write cdr
	s.CDRLog.InfoContext(ctx, "cdr", slog.Any("cdr", nc))

	// step 6 save
	_, err = s.DBx.SaveCurrent(ctx, nc)
	return nc, err
}
