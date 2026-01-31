package hank

import (
	"encoding/json/v2"
	"net/http"
)

const (
	deviceList  = "deviceList"
	gatewayInfo = "gatewayInfo"
	rate        = "rate"
	deviceData  = "deviceData"
)

func Handle(h *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data SyncData
		if err := json.UnmarshalRead(r.Body, &data); err != nil {
			return
		}

		ctx := r.Context()

		var err error

		switch data.Type {
		case deviceList:
			err = h.HandleSyncDeviceInfo(ctx, data)
		case gatewayInfo:
			err = h.HandleUploadGatewayInfo(ctx, data)
		case deviceData:
			err = h.HandleSyncDeviceData(ctx, data)
		case rate:
			err = ErrNoRate
		default:
			err = Error(data.Type)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.MarshalWrite(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.MarshalWrite(w, OK)
	}
}
