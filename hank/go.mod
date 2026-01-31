module github.com/twiglab/h2o/hank

go 1.25.5

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/google/uuid v1.6.0
	github.com/twiglab/h2o v0.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
)

replace github.com/twiglab/h2o => ../
