package servicemonitor

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventEmitter interface {
	Emit(event string, data interface{})
}

type WailsEmitter struct {
	ctx context.Context
}

func NewEventEmitter(ctx context.Context) *WailsEmitter {
	return &WailsEmitter{ctx: ctx}
}

func (e *WailsEmitter) Emit(event string, data interface{}) {
	wailsRuntime.EventsEmit(e.ctx, event, data)
}
