package idb

import (
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

const fname = "filename"

func parseCSVArgs(args []string) (file string) {
	for _, s := range args {
		ss := strings.Split(s, "=")
		if len(ss) == 2 {
			switch ss[0] {
			case fname:
				return unquote(ss[1])
			}
		}
	}
	return ""
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func loadTable[T VitrualRecord](file string) (*VirtualTable[T], error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rs []T

	if err := gocsv.UnmarshalFile(f, &rs); err != nil {
		return nil, err
	}

	t := &VirtualTable[T]{
		Rows: rs,
	}
	return t, nil
}
