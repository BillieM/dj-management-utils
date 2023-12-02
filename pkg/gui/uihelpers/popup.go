package uihelpers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewPercentagePopup(title string, content fyne.CanvasObject, win fyne.Window, resizeEvents *ResizeEvents, percentW, percentH float32, maxSize fyne.Size, onClose func(close func())) *widget.PopUp {

	var popup *widget.PopUp

	key := resizeEvents.Add(func() {})

	header := container.NewHBox()

	header.Add(&widget.Label{Text: title, TextStyle: fyne.TextStyle{Bold: true}})
	header.Add(layout.NewSpacer())
	header.Add(&widget.Button{Icon: theme.ContentClearIcon(), Importance: widget.LowImportance, OnTapped: func() {
		if onClose != nil {
			onClose(popup.Hide)
		}
		resizeEvents.Remove(key)
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

	resizeEvents.Set(key, func() {
		popup.Resize(CanvasPercentSize(win, percentW, percentH, popupContent.MinSize(), maxSize))
	})

	popup.Resize(CanvasPercentSize(win, percentW, percentH, popupContent.MinSize(), maxSize))

	return popup
}
