package uihelpers

import "fyne.io/fyne/v2"

type ResizeEvents struct {
	parentWindow    fyne.Window
	resizeCallbacks []*func()
}

func NewResizeEvents(parentWindow fyne.Window, resizeCallbacks []*func()) *PercentagePopupLayout {
	return &PercentagePopupLayout{
		parentWindow:    parentWindow,
		resizeCallbacks: resizeCallbacks,
	}
}

func (p *ResizeEvents) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (p *ResizeEvents) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, callback := range p.resizeCallbacks {
		(*callback)()
	}
}
