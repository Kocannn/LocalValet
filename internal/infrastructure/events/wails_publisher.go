package events

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type WailsPublisher struct {
	ctx context.Context
}

func NewWailsPublisher(ctx context.Context) *WailsPublisher {
	return &WailsPublisher{ctx: ctx}
}

func (p *WailsPublisher) Emit(event string, data interface{}) {
	wailsRuntime.EventsEmit(p.ctx, event, data)
}
