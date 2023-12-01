package iwidget

import (
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
Opens a url in the default browser.

If the 'SEREN_USE_CHROME flag is set, then the url will be opened in Chrome.
This flag is not set by default, and is only used for testing purposes.
*/
type OpenInBrowserButton struct {
	widget.Button

	URL *url.URL
}

func NewOpenInBrowserButton(text string, urlString string) *OpenInBrowserButton {

	openInBrowserBtn := &OpenInBrowserButton{}

	if urlString != "" {
		openInBrowserBtn.SetContent(text, urlString)
	}

	openInBrowserBtn.ExtendBaseWidget(openInBrowserBtn)

	return openInBrowserBtn
}

func (i *OpenInBrowserButton) setOpenFunc() {
	if os.Getenv("SEREN_USE_CHROME") == "" {
		i.OnTapped = func() {
			err := fyne.CurrentApp().OpenURL(i.URL)
			if err != nil {
				fyne.LogError("Failed to open url", err)
			}
		}
	} else {
		i.OnTapped = func() {
			helpers.OpenInChrome(i.URL)
		}
	}
}

func (i *OpenInBrowserButton) SetContent(text, urlStr string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Failed to parse url", err)
	}
	i.URL = u
	i.setOpenFunc()
	i.SetText(text)
}

type ClickableLabel struct {
	widget.Label

	OnTapped func()
}

func NewClickableLabel(text string, tapped func()) *ClickableLabel {

	clickableLabel := &ClickableLabel{
		OnTapped: tapped,
	}

	clickableLabel.ExtendBaseWidget(clickableLabel)

	return clickableLabel
}

func (i *ClickableLabel) Tapped() {
	i.OnTapped()
}
