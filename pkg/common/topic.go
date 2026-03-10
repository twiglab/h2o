package common

import "strings"

const H2O = "h2o"

const (
	ELECTRICITY = "E"
	WATER       = "W"
	GAS         = "G"
)

const (
	WaterTopic       = "h2o/+/W"
	ElectricityTopic = "h2o/+/E"
	GasTopic         = "h2o/+/G"
)

func Topic(d Device) string {
	return H2O + "/" + d.Code + "/" + d.Type
}

func TopicPart(topic string) (string, string, string) {
	parts := strings.SplitN(topic, "/", 3)
	_ = parts[2]
	return parts[0], parts[1], parts[2]
}

func TopicType(topic string) string {
	_, _, t := TopicPart(topic)
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
