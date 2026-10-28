package vigil

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/clog/wal"
)

type Hub struct {
	// TSDB Recorder
	DB Recorder

	Logger *slog.Logger

	BaseContext func(h *Hub) context.Context

	WAL *wal.WAL
}

func (h *Hub) Handle(ctx context.Context, data Meter) error {

	h.WAL.WriteLogContext(ctx,
		wal.String("type", data.Type),
		wal.Any("data", data),
	)

	// h.TSDB.Tabb(ctx, data.Code, data)

	if err := h.DB.Tabb(ctx, data.Code, data); err != nil {
		h.Logger.ErrorContext(ctx, "Record Electy error", slog.Any("data", data), slog.Any("error", err))
		return err
	}
	return nil
}
