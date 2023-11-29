package iwidget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func SimpleDialog(title string, icon fyne.Resource, content fyne.CanvasObject, win fyne.Window, onClose func(close func())) *widget.PopUp {
	var popup *widget.PopUp

	header := container.NewHBox()
	if icon != nil {
		header.Add(widget.NewIcon(icon))
	}
	header.Add(&widget.Label{Text: title, TextStyle: fyne.TextStyle{Bold: true}})
	header.Add(layout.NewSpacer())
	header.Add(&widget.Button{Icon: theme.ContentClearIcon(), Importance: widget.LowImportance, OnTapped: func() {
		if onClose == nil {
			popup.Hide()
		} else {
			onClose(popup.Hide)
		}
	}})

	popup = widget.NewModalPopUp(container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		nil, nil, nil,
		content,
	), win.Canvas())

	return popup
}
