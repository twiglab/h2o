package s7

import (
	"github.com/twiglab/h2o/nab/box/driverbox"
	"github.com/twiglab/h2o/nab/box/plugins/s7/internal"
)

func EnablePlugin() {
	driverbox.EnablePlugin(internal.ProtocolName, new(internal.Plugin))
}
