package box

import (
	"github.com/simonvetter/modbus"
)

func NewModbusClient(config *modbus.ClientConfiguration) (*modbus.ModbusClient, error) {
	var client *modbus.ModbusClient
	var err error

	if client, err = modbus.NewClient(config); err != nil {
		return client, err
	}

	return client, client.Open()
}
