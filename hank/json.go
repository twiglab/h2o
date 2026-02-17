package hank

import (
	"encoding/json/v2"
	"io"
)

func unmarshal(in []byte, out any) error {
	return json.Unmarshal(in, out)
}

func marshalWrite(out io.Writer, in any) error {
	return json.MarshalWrite(out, in)
}

func marshal(in any) ([]byte, error) {
	return json.Marshal(in)
}
