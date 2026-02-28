package chrgg

import (
	"context"

	"github.com/twiglab/h2o/wal"
)

type ChargeServer struct {
	DBx         *DBx
	CdrWAL      *wal.WAL
	ChargEngine ChargeEngine
	CheckFunc   CheckFunc
	VerifyFunc  VerifyFunc
}

func (s *ChargeServer) pre(_ context.Context, md MeterData) (ChargeData, error) {
	return ChargeData{MeterData: md}, nil
}

func (s *ChargeServer) Charge(ctx context.Context, md MeterData) (CDR, error) {
	// setp 1 prepare
	cd, err := s.pre(ctx, md)
	if err != nil {
		return nilCDR, err
	}

	// step 2 load
	l, _, err := s.DBx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return nilCDR, err
	}
	last := MakeLast(l)

	// step 3 verify and check
	if !s.VerifyFunc(ctx, last, cd) {
		return nilCDR, nil
	}

	if err := s.CheckFunc(ctx, last, cd); err != nil {
		return nilCDR, err
	}

	// setp 4 calc
	ru, err := s.ChargEngine.GetRuler(ctx, cd)
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
