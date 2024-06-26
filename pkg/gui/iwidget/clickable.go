package iwidget

import (
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
Contains a selection of clickable widgets, such as buttons and labels.
*/

/*
Opens a url in the default browser.

If the 'SEREN_USE_CHROME flag is set, then the url will be opened in Chrome.
This flag is not set by default, and is only used for testing purposes.
*/
type OpenInBrowserButton struct {
	*Base
	widget.Button

	URL *url.URL
}

func NewOpenInBrowserButton(widgetBase *Base, text string, urlString string) *OpenInBrowserButton {

	openInBrowserBtn := &OpenInBrowserButton{
		Base: widgetBase,
	}

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

func (i *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	i.OnTapped()
}

type ClickableCard struct {
	widget.Card

	OnTapped func()
}

func NewClickableCard(tapped func()) *ClickableCard {

	clickableCard := &ClickableCard{
		OnTapped: tapped,
	}

	clickableCard.ExtendBaseWidget(clickableCard)

	return clickableCard
}

func (i *ClickableCard) Tapped(_ *fyne.PointEvent) {
	i.OnTapped()
}
