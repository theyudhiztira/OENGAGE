package broadcast

import (
	"context"
)

type broadcastHandler struct {
	Ctx     *context.Context
	Service broadcastService
}

func NewBroadcastHandler(ctx *context.Context, service broadcastService) *broadcastHandler {
	return &broadcastHandler{
		Ctx:     ctx,
		Service: service,
	}
}
