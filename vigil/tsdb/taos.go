package tsdb

import (
	"unsafe"
)

const (
	TAG_CODE  = "code"
	TAG_PROJ  = "proj"
	TAG_PCODE = "pcode"

	FIELD_B = "b"

	FIELD_DATA_VALUE = "dv" // 表显

	FIELD_FREQUENCY = "f"

	FIELD_I_A = "ia"
	FIELD_I_B = "ib"
	FIELD_I_C = "ic"

	FIELD_P = "p" // 总有功功率

	FIELD_V_A = "va"
	FIELD_V_B = "vb"
	FIELD_V_C = "vc"
)

const (
	ELECTY_STB = "electy_stb"
	WATER_STB  = "water_stb"
)

const (
	TSDB_SML_TIMESTAMP_SECONDS       = "s"
	TSDB_SML_TIMESTAMP_MILLI_SECONDS = "ms"
	TSDB_SML_TIMESTAMP_MICRO_SECONDS = "us"
	TSDB_SML_TIMESTAMP_NANO_SECONDS  = "ns"
)

func bytesToStr(bs []byte) string {
	return unsafe.String(&bs[0], len(bs))
}
