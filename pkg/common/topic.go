package common

func Topic(d Device) string {
	return "h2o/" + d.Code + "/" + d.Type
}

const (
	WaterTopic       = "h2o/+/W"
	ElectricityTopic = "h2o/+/E"
	GasTopic         = "h2o/+/G"
)
