package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

type ChargeData struct {
	RawData
}

type ChangeServer struct {
	dbx    *DBx
	cdrLog *slog.Logger

	re RulerEngine
}

func (s *ChangeServer) pre(_ context.Context, rd RawData) (ChargeData, error) {
	return ChargeData{RawData: rd}, nil
}

func (s ChangeServer) check(ctx context.Context, last *ent.CDR, cd ChargeData) error {
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

func (s *ChangeServer) DoChange(ctx context.Context, bd RawData) (nc CDR, err error) {
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
