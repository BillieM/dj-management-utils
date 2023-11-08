package ui

import "fyne.io/fyne/v2/widget"

type MyProgressBar struct {
	*widget.ProgressBar
}

func (d *Data) buildProgressBar() MyProgressBar {
	return MyProgressBar{
		widget.NewProgressBar(),
	}
}

func (p MyProgressBar) updateProgressBar(value float64) {
	p.SetValue(value)
}
