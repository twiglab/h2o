package main

import (
	"github.com/twiglab/h2o/nab/cmd/nab/cmd"
	"github.com/twiglab/h2o/nab/equlib"
)

func main() {
	equlib.MustHas()
	cmd.Execute()
}
