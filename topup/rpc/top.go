package rpc

import (
	"context"
	"encoding/json"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
	"golang.org/x/exp/jsonrpc2"
)

type topHandle struct {
	cli *ent.Client
}

func (h topHandle) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	var p PosReq
	if err := json.Unmarshal(req.Params, &p); err != nil {
		return nil, err
	}

	q := h.cli.Device.Query()
	if p.PosCode != "" {
		q.Where(device.PosCodeEQ(p.PosCode))
	}
	ds, err := q.All(ctx)

	return &PosResp{Data: ds}, err
}
