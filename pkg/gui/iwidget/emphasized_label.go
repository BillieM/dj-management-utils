package iwidget

import "fyne.io/fyne/v2/widget"

type EmphasizedLabel struct {
	widget.Label
}

func NewEmphasizedLabel(text string) *EmphasizedLabel {
	i := &EmphasizedLabel{}
	i.TextStyle.Bold = true
	i.Importance = widget.HighImportance

	i.ExtendBaseWidget(i)

	return i
}
