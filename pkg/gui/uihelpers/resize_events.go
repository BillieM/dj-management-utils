package uihelpers

import (
	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type ResizeEvents struct {
	resizeCallbacks map[string]func()
}

func NewResizeEvents() *ResizeEvents {
	return &ResizeEvents{
		resizeCallbacks: make(map[string]func()),
	}
}

func (p *ResizeEvents) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (p *ResizeEvents) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, cb := range p.resizeCallbacks {
		cb()
	}
}

func (p *ResizeEvents) Add(callback func()) string {
	key := uuid.New().String()

	p.resizeCallbacks[key] = callback

	return key
}

func (p *ResizeEvents) Remove(key string) {
	delete(p.resizeCallbacks, key)
}

func (p *ResizeEvents) Clear() {
	p.resizeCallbacks = make(map[string]func())
}

func (p *ResizeEvents) Set(key string, callback func()) {
	p.resizeCallbacks[key] = callback
}
