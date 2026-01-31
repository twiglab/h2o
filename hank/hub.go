package hank

import (
	"context"
	"encoding/json/v2"
	"log/slog"

	"github.com/twiglab/h2o/pkg/kwh"
)

type Enh struct {
}

func (e *Enh) Convert(dd DeviceData) kwh.Device {
	return kwh.Device{
		Code: dd.No,
		Type: dd.Type,
		Name: dd.No,

		Time: dd.LastDataTime,

		Data: kwh.Data{},
	}
}

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
		h.DataLog.DebugContext(ctx, "deviceData", slog.Any("data", dd))
		kwhd := h.Enh.Convert(dd)
		_ = h.Sender.SendData(ctx, kwhd)
	}

	return nil
}
