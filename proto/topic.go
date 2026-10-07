package proto

import "strings"

const H2O = "h2o"

const (
	DATA = "data"
	OPT  = "opt"

	ONOFF = "onoff"
	QUOTA = "quota"
)

const (
	ON  = "ON"
	OFF = "OFF"
)

const (
	TYPE_ELECTRICITY = "E"
	TYPE_WATER       = "W"
	TYPE_GAS         = "G"
)

func TopicPart(topic string) []string {
	parts := strings.Split(topic, "/")
	_ = parts[1]

	if parts[0] != H2O {
		panic("not h2o topic")
	}
	return parts
}
