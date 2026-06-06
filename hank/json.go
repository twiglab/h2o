package hank

import (
	"encoding/json/v2"
	"io"
)

func unmarshal(in []byte, out any) error {
	return json.Unmarshal(in, out)
}

func marshal(in any) ([]byte, error) {
	return json.Marshal(in)
}

var newline = []byte{'\n'}

func writeReturn(out io.Writer, in any) (err error) {
	if err = json.MarshalWrite(out, in); err == nil {
		_, err = out.Write(newline)
	}
	return
}
