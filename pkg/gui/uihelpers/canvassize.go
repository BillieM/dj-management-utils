package uihelpers

import (
	"fyne.io/fyne/v2"
)

/*
CanvasPercentSize returns a fyne.Size struct with the width and height set to
the percentage of the canvas size passed in

Thank you to https://github.com/matwachich for this function <3
*/
func CanvasPercentSize(win fyne.Window, percentW, percentH float32, minSize fyne.Size, maxSize fyne.Size) fyne.Size {
	csz := win.Canvas().Size()

	size := fyne.NewSize(0, 0)

	if percentW > 0 {
		if percentW > 1 {
			percentW /= 100
		}

		w := csz.Width * percentW

		if maxSize.Width > 0 && w > maxSize.Width {
			size.Width = maxSize.Width
		} else if minSize.Width > 0 && w < minSize.Width {
			size.Width = minSize.Width
		} else {
			size.Width = w
		}
	}

	if percentH > 0 {
		if percentH > 1 {
			percentH /= 100
		}

		h := csz.Height * percentH

		if maxSize.Height > 0 && h > maxSize.Height {
			size.Height = maxSize.Height
		} else if minSize.Height > 0 && h < minSize.Height {
			size.Height = minSize.Height
		} else {
			size.Height = h
		}
	}

	return size
}
