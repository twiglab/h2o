package hank

import (
	"context"
	"encoding/json/v2"
	"log/slog"
)

type Hub struct {
	DataLog *slog.Logger
	InfoLog *slog.Logger
	Enh     *Enh
	Sender  *MQTTAction
}

func (h *Hub) HandleSyncDeviceInfo(ctx context.Context, data SyncData) error {
	var dl DeviceList
	if err := json.Unmarshal(data.Data, &dl); err != nil {
		return err
	}
	for _, di := range dl {
		h.InfoLog.DebugContext(ctx, "deviceInfo", slog.Any("data", di))
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
		kwhd := h.Enh.Convert(dd)
		h.DataLog.DebugContext(ctx, "deviceData", slog.Any("data", kwhd))
		if h.Sender != nil {
			_ = h.Sender.SendData(ctx, kwhd)
		}
	}

	return nil
}
