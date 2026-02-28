package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/wal"
)

type ChargeServer struct {
	DBx         *DBx
	CdrWAL      *wal.WAL
	ChargEngine ChargeEngine
	CheckFunc   CheckFunc
	VerifyFunc  VerifyFunc

	Logger *slog.Logger
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
		s.Logger.ErrorContext(ctx, "loadLast error", slog.Any("chargeData", cd), slog.Any("error", err))
		return nilCDR, err
	}
	last := MakeLast(l)

	// step 3 verify and check
	if vr, ok := s.VerifyFunc(ctx, last, cd); !ok {
		s.Logger.DebugContext(ctx, "verify", slog.Any("last", last), slog.Any("cd", cd), slog.Any("return", vr))
		return nilCDR, nil
	}

	if err := s.CheckFunc(ctx, last, cd); err != nil {
		s.Logger.ErrorContext(ctx, "check error", slog.Any("last", last), slog.Any("chargeData", cd), slog.Any("error", err))
		return nilCDR, err
	}

	// setp 4 calc
	ru, err := s.ChargEngine.GetRuler(ctx, cd)
	if err != nil {
		s.Logger.ErrorContext(ctx, "GerRuler error", slog.Any("error", err), slog.Any("cd", cd))
		return nilCDR, err
	}

	nc := CalcCDR(last, cd, ru)

	// step 5 write cdr
	s.CdrWAL.WriteLogContext(ctx, wal.Type("nhcdr"), wal.Data(nc))

	// step 6 save
	_, err = s.DBx.SaveCurrent(ctx, nc)
	if err != nil {
		s.Logger.ErrorContext(ctx, "save error", slog.Any("ncdr", nc), slog.Any("error", err))
	}
	return nc, err
}
