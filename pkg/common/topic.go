package common

import "strings"

func Topic(d Device) string {
	return "h2o/" + d.Code + "/" + d.Type
}

const (
	WaterTopic       = "h2o/+/W"
	ElectricityTopic = "h2o/+/E"
	GasTopic         = "h2o/+/G"
)

func TopicType(topic string) string {
	t, _ := strings.CutSuffix(topic, "/")
	switch t {
	case WATER:
		return WaterTopic
	case ELECTRICITY:
		return ElectricityTopic
	case GAS:
		return GasTopic
	}
	panic(topic + " not supports")
}
