package main

import (
	"github.com/twiglab/h2o/nab/equlib"
	_ "github.com/twiglab/h2o/nab/equlib/clfoc"
	_ "github.com/twiglab/h2o/nab/equlib/kh"
)

func init() {
	equlib.MustHas()
}
