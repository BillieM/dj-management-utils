package uihelpers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewPercentagePopup(title string, content fyne.CanvasObject, win fyne.Window, contentStack *fyne.Container, percentW, percentH float32, onClose func(close func())) *widget.PopUp {

	resizeLayout := NewPercentagePopupLayout(
		win,
		func() {},
		0.9, 0.9,
	)

	resizeContainer := container.New(resizeLayout, container.NewStack())

	var popup *widget.PopUp

	header := container.NewHBox()

	header.Add(&widget.Label{Text: title, TextStyle: fyne.TextStyle{Bold: true}})
	header.Add(layout.NewSpacer())
	header.Add(&widget.Button{Icon: theme.ContentClearIcon(), Importance: widget.LowImportance, OnTapped: func() {
		if onClose != nil {
			onClose(popup.Hide)
		}
		contentStack.Remove(resizeContainer)
		popup.Hide()
	}})

	popupContent := container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		nil, nil, nil,
		content,
	)

	popup = widget.NewModalPopUp(
		popupContent,
		win.Canvas(),
	)

	contentStack.Add(resizeContainer)

	resizeLayout.popupCallback = func() {
		popup.Resize(CanvasPercentSize(win, percentW, percentH, popup.MinSize()))
	}

	return popup
}

type PercentagePopupLayout struct {
	parentWindow       fyne.Window
	popupCallback      func()
	percentW, percentH float32
}

func NewPercentagePopupLayout(parentWindow fyne.Window, popupCallback func(), percentW, percentH float32) *PercentagePopupLayout {
	return &PercentagePopupLayout{
		parentWindow:  parentWindow,
		percentW:      percentW,
		percentH:      percentH,
		popupCallback: popupCallback,
	}
}

func (p *PercentagePopupLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (p *PercentagePopupLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	p.popupCallback()
}
