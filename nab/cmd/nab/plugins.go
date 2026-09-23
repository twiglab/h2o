package main

import (
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/equlib/fake"

	_ "github.com/twiglab/h2o/nab/equlib/clfoc"
	_ "github.com/twiglab/h2o/nab/equlib/kh"
)

func init() {
	equlib.Register("fake-device", fake.New(false))
	equlib.Register("fake-device-w", fake.New(true))
}
