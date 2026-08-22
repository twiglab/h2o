package box

import (
	"context"
	"log/slog"
)

type LogAction struct {
}

func (c LogAction) SendData(ctx context.Context, obj SendObject) error {
	slog.DebugContext(ctx, "logAction", slog.Any("data", obj), slog.String("topic", obj.Topic()))
	return nil
}
