package helpers

import "fyne.io/fyne/v2/widget"

type Importance int

const (
	MediumImportance Importance = iota
	LowImportance
	HighImportance
	WarningImportance
	DangerImportance
	SuccessImportance
)

func (i Importance) String() string {
	switch i {
	case LowImportance:
		return "LowImportance"
	case MediumImportance:
		return "MediumImportance"
	case HighImportance:
		return "HighImportance"
	case WarningImportance:
		return "WarningImportance"
	case DangerImportance:
		return "DangerImportance"
	case SuccessImportance:
		return "SuccessImportance"
	default:
		return "Unknown Importance"
	}
}

func (i Importance) GetFyneImportance() widget.Importance {
	switch i {
	case LowImportance:
		return widget.LowImportance
	case MediumImportance:
		return widget.MediumImportance
	case HighImportance:
		return widget.HighImportance
	case WarningImportance:
		return widget.WarningImportance
	case DangerImportance:
		return widget.DangerImportance
	case SuccessImportance:
		return widget.SuccessImportance
	default:
		return widget.MediumImportance
	}
}
