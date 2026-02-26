package hank

import (
	"bytes"
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

func writeReturn(out io.Writer, in any) error {
	var bf bytes.Buffer
	bf.Grow(256)

	if err := json.MarshalWrite(&bf, in); err != nil {
		return err
	}

	if err := bf.WriteByte('\n'); err != nil {
		return err
	}

	_, err := bf.WriteTo(out)

	return err
}
