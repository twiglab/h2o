package chrgg

import (
	"context"
	"log/slog"
)

type ChargeData struct {
	OrgiData
}

type ChangeServer struct {
	dbx    *DBx
	cdrLog *slog.Logger

	re RulerEngine
}

func (s *ChangeServer) pre(ctx context.Context, bcd OrgiData) (ChargeData, error) {
	return ChargeData{OrgiData: bcd}, nil
}

func (s *ChangeServer) calc(ctx context.Context, cd ChargeData) (CDR, error) {
	last, notfound, err := s.dbx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return Nil, err
	}

	if notfound {
		return FirstCDR(cd), nil
	}

	ru, err := s.re.GetResult(ctx, cd)
	if err != nil {
		return Nil, err
	}

	return CalcCDR(last, cd, ru), nil
}

func (s *ChangeServer) DoChange(ctx context.Context, bd OrgiData) (nc CDR, err error) {
	var cd ChargeData
	if cd, err = s.pre(ctx, bd); err != nil {
		return
	}
	if nc, err = s.calc(ctx, cd); err == nil { // err == nil
		_, err = s.dbx.SaveCurrent(ctx, nc)
	}
	return
}
