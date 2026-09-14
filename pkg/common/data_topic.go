package common

import "strings"

const H2O = "h2o"

const (
	TYPE_ELECTRICITY = "E"
	TYPE_WATER       = "W"
	TYPE_GAS         = "G"
)

const (
	WaterDataTopic       = "h2o/data/+/W"
	ElectricityDataTopic = "h2o/data/+/E"
	GasDataTopic         = "h2o/data/+/G"
)

func DataTopic(d Device) string {
	return H2O + "/data/" + d.Code + "/" + d.Type
}

func TopicPart(topic string) []string {
	parts := strings.Split(topic, "/")
	return parts

}

func DataTopicPart(topic string) (string, string, string, string) {
	parts := TopicPart(topic)
	_ = parts[3]
	return parts[0], parts[1], parts[2], parts[3]
}

func DataTopicType(topic string) string {
	_, _, _, t := DataTopicPart(topic)
	switch t {
	case TYPE_WATER:
		return WaterDataTopic
	case TYPE_ELECTRICITY:
		return ElectricityDataTopic
	case TYPE_GAS:
		return GasDataTopic
	}
	panic(topic + " not supports")
}
