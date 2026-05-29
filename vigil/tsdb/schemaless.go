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

func (s *Schemaless) TabbElecty(ctx context.Context, data vigil.ElectricityMeter) error {

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Second)

	enc.StartLine(ELECTY_STB)

	enc.AddTag(TAG_CODE, data.Code)
	enc.AddTag(TAG_PCODE, data.Pos.PosCode)
	enc.AddTag(TAG_PROJ, data.Pos.Project)

	v, _ := lineprotocol.FloatValue(data.STD)
	enc.AddField(FIELD_B, v)

	enc.AddField(FIELD_DATA_VALUE, lineprotocol.IntValue(data.Data.DataValue))

	enc.AddField(FIELD_FREQUENCY, lineprotocol.IntValue(data.Data.Frequency))

	enc.AddField(FIELD_I_A, lineprotocol.IntValue(data.Data.CurrentA))
	enc.AddField(FIELD_I_B, lineprotocol.IntValue(data.Data.CurrentB))
	enc.AddField(FIELD_I_C, lineprotocol.IntValue(data.Data.CurrentC))

	enc.AddField(FIELD_P, lineprotocol.IntValue(data.Data.ActivePowerTotal))

	enc.AddField(FIELD_V_A, lineprotocol.IntValue(data.Data.VoltageA))
	enc.AddField(FIELD_V_B, lineprotocol.IntValue(data.Data.VoltageB))
	enc.AddField(FIELD_V_C, lineprotocol.IntValue(data.Data.VoltageC))

	enc.EndLine(data.DataTime)

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.SchemalessInsert(common.GetReqID(), line, unified.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_SECONDS, 0, "")
}

func (s *Schemaless) TabbWater(ctx context.Context, data vigil.WaterMeter) error {
	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Second)

	enc.StartLine(WATER_STB)

	enc.AddTag(TAG_CODE, data.Code)
	enc.AddTag(TAG_PCODE, data.Pos.PosCode)
	enc.AddTag(TAG_PROJ, data.Pos.Project)

	enc.AddField(FIELD_DATA_VALUE, lineprotocol.IntValue(data.Data.DataValue))

	enc.EndLine(data.DataTime)

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.SchemalessInsert(common.GetReqID(), line, unified.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_SECONDS, 0, "")
}
