package plugins

import (
	"github.com/twiglab/h2o/nab/box/plugins/httpclient"
	"github.com/twiglab/h2o/nab/box/plugins/httpserver"
	"github.com/twiglab/h2o/nab/box/plugins/modbus"
	"github.com/twiglab/h2o/nab/box/plugins/mqtt"
	"github.com/twiglab/h2o/nab/box/plugins/s7"
	"github.com/twiglab/h2o/nab/box/plugins/tcpserver"
	"github.com/twiglab/h2o/nab/box/plugins/websocket"
)

func EnableAll() {
	modbus.EnablePlugin()
	//bacnet.EnablePlugin()
	httpserver.EnablePlugin()
	httpclient.EnablePlugin()
	websocket.EnablePlugin()
	tcpserver.EnablePlugin()
	mqtt.EnablePlugin()
	//dlt645.EnablePlugin()
	//opcua.EnablePlugin()
	s7.EnablePlugin()
}
