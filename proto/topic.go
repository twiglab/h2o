package proto

import (
	"fmt"
	"strings"
)

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

func TopicPart(topic string) ([]string, error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 2 || parts[0] != H2O {
		return nil, fmt.Errorf("invalid h2o topic: %q", topic)
	}
	return parts, nil
}

func MakeTopic(typ string, parts ...string) string {
	c := append([]string{H2O, typ}, parts...)
	return strings.Join(c, "/")
}
