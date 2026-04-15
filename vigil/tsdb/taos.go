package tsdb

import (
	"unsafe"
)

const (
	TAG_CODE  = "code"
	TAG_PROJ  = "proj"
	TAG_PCODE = "pcode"

	FIELD_DENSITY_COUNT = "human_count"
	FIELD_DENSITY_RATIO = "human_ratio"
)

const (
	POWER_STB = "power_stb"
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
