package hank

import (
	"context"
	"encoding/json/v2"
	"log/slog"
)

type Enhancer interface {
}

type Hub struct {
	dataLog *slog.Logger
	infoLog *slog.Logger
	enhance Enhancer
}

func (h *Hub) HandleSyncDeviceInfo(ctx context.Context, data SyncData) error {
	var dl DeviceList
	if err := json.Unmarshal(data.Data, &dl); err != nil {
		return err
	}
	for _, di := range dl {
		h.infoLog.DebugContext(ctx, "deviceInfo", slog.Any("data", di))
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

	for _, dd := range ddl {
		h.dataLog.DebugContext(ctx, "deviceData", slog.Any("data", dd))
	}

	return nil
}
