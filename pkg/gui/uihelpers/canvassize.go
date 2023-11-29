package uihelpers

import "fyne.io/fyne/v2"

/*
CanvasPercentSize returns a fyne.Size struct with the width and height set to
the percentage of the canvas size passed in

Thank you to https://github.com/matwachich for this function <3
*/
func CanvasPercentSize(w fyne.Window, percentW, percentH float32, minSize fyne.Size) fyne.Size {
	csz := w.Canvas().Size()

	if percentW > 0 {
		if percentW > 1 {
			percentW /= 100
		}

		w := csz.Width * percentW
		if w > minSize.Width {
			minSize.Width = w
		}
	}

	if percentH > 0 {
		if percentH > 1 {
			percentH /= 100
		}

		h := csz.Height * percentH
		if h > minSize.Height {
			minSize.Height = h
		}
	}

	return minSize
}
