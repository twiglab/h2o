package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/wal"
)

type ChargeServer struct {
	DBx    *DBx
	CdrWAL *wal.WAL
	Eng    ChargeEngine
}

func (s *ChargeServer) pre(_ context.Context, md MeterData) (ChargeData, error) {
	return ChargeData{MeterData: md}, nil
}

func (s ChargeServer) check(_ context.Context, last *ent.CDR, cd ChargeData) error {
	err := checkDup(last, cd)
	return err
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

func (s *ChargeServer) doNewCharge(ctx context.Context, cd ChargeData) (CDR, error) {
	nc := CalcCDR(first, cd, RulNew)
	s.CdrWAL.WriteLogContext(ctx, wal.Type("nhcdr"), wal.Data(cd))
	_, err := s.DBx.SaveCurrent(ctx, nc)
	return nc, err
}

func (s *ChargeServer) doCharge(ctx context.Context, last *ent.CDR, cd ChargeData) (CDR, error) {
	// step 3 verify and check
	if !s.Verify(ctx, last, cd) {
		return nilCDR, nil
	}

	if err := s.check(ctx, last, cd); err != nil {
		return nilCDR, err
	}

	// setp 4 calc
	ru, err := s.Eng.GetRuler(ctx, cd)
	if err != nil {
		return nilCDR, err
	}

	nc := CalcCDR(last, cd, ru)

	// step 5 write cdr
	s.CdrWAL.WriteLogContext(ctx, wal.Type("nhcdr"), wal.Data(cd))

	// step 6 save
	_, err = s.DBx.SaveCurrent(ctx, nc)
	return nc, err
}

func (s *ChargeServer) DoCharge(ctx context.Context, bd MeterData) (CDR, error) {
	// setp 1 prepare
	cd, err := s.pre(ctx, bd)
	if err != nil {
		return nilCDR, err
	}

	// step 2 load
	last, notfound, err := s.DBx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return nilCDR, err
	}

	// step 2.1 save first
	if notfound {
		return s.doNewCharge(ctx, cd)
	}

	// doCharge
	return s.doCharge(ctx, last, cd)
}
