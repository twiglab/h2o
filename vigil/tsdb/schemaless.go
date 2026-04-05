package tsdb

import (
	"context"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"github.com/taosdata/driver-go/v3/common"
	"github.com/taosdata/driver-go/v3/ws/unified"
	"github.com/twiglab/h2o/vigil"
)

type Schemaless struct {
	schemaless *unified.Client
}

func NewSchLe(dsn string) (*Schemaless, error) {
	s, err := unified.Open(dsn)
	if err != nil {
		return nil, err
	}

	return &Schemaless{
		schemaless: s,
	}, nil
}

func (s *Schemaless) HandleElecty(ctx context.Context, data vigil.ElectricityMeter) error {

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Second)
	enc.StartLine(POWER_STB)

	enc.AddTag(TAG_CODE, data.Code)
	enc.AddTag(TAG_PCODE, data.Pos.PosCode)
	enc.AddTag(TAG_PROJ, data.Pos.Project)

	enc.EndLine(data.DataTime)

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.SchemalessInsert(common.GetReqID(), line, unified.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_SECONDS, 0, "")
}
