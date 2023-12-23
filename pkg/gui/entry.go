package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
)

func Entry() {
	a := app.New()
	mainWindow := a.NewWindow("Library Utilities")

	e, err := buildGuiEnv(a, mainWindow)

	if err != nil {
		log.Fatal(
			fault.Wrap(err, fmsg.With("Error building GUI environment, exiting")),
		)
	}

	appLoading := iwidget.NewAppLoading()

	mainWindow.SetContent(
		container.NewStack(
			appLoading,
		),
	)

	go e.loadApp(func() {
		// Seems strange this method is called SetMainMenu as it really defines the top bar of the application, but hey :)
		mainWindow.SetMainMenu(e.makeNavBar(a, mainWindow))

		mainWindow.Resize(fyne.NewSize(1600, 900))

		contentStack := container.NewStack()

		e.setMainContent(contentStack, e.getViewList()["home"])

		split := container.NewHSplit(e.makeNavMenu(contentStack), contentStack)
		split.SetOffset(0)

		mainWindow.SetContent(
			container.NewStack(
				container.New(e.resizeEvents, container.NewStack()),
				split,
			),
		)
	})

	mainWindow.ShowAndRun()
}
