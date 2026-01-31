package hank

import (
	"context"
	"encoding/json/v2"
)

type Hub struct {
}

func (h *Hub) HandleSyncDeviceInfo(ctx context.Context, data SyncData) error {
	var dl DeviceList
	if err := json.Unmarshal(data.Data, &dl); err != nil {
		return err
	}
	return nil
}

func (h *Hub) HandleUploadGatewayInfo(ctx context.Context, data SyncData) error {
	return nil
}

func (h *Hub) HandleSyncDeviceData(ctx context.Context, data SyncData) error {
	var ddl DeviceDataList
	if err := json.Unmarshal(data.Data, &ddl); err != nil {
		return err
	}
	return nil
}
