package websocket

import (
	"github.com/twiglab/h2o/nab/box/driverbox"
	"github.com/twiglab/h2o/nab/box/plugins/websocket/internal"
)

func EnablePlugin() {
	driverbox.EnablePlugin(internal.ProtocolName, new(internal.Plugin))
}
