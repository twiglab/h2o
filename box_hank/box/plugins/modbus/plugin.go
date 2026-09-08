package modbus

import (
	"github.com/twiglab/h2o/nab/box/driverbox"
	"github.com/twiglab/h2o/nab/box/plugins/modbus/internal"
)

func EnablePlugin() {
	driverbox.EnablePlugin(internal.ProtocolName, new(internal.Plugin))
}
