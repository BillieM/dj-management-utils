package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/helpers"
)

/*
makeNavMenu builds the navigation menu on the left side of the application

the navigation menu is a tree object and is dynamically built from the operations list inside data.go
*/

func (d *Data) makeNavMenu(w fyne.Window, contentStack *fyne.Container) fyne.CanvasObject {

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return d.OperationIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := d.OperationIndex[uid]
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Node")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			op, ok := d.Operations[uid]
			if !ok {
				helpers.HandleFatalError(helpers.ErrOperationNotFound)
				return
			}
			node.(*widget.Label).SetText(op.Name)
		},
		OnSelected: func(uid string) {
			op, ok := d.Operations[uid]
			if !ok {
				helpers.HandleFatalError(helpers.ErrOperationNotFound)
				return
			}
			if d.processing {
				showErrorDialog(w, helpers.ErrPleaseWaitForProcess)
				return
			}
			d.setMainContent(w, contentStack, op)
		},
	}

	tree.OpenAllBranches()

	navContainer := container.NewBorder(widget.NewLabel("menu <3"), nil, nil, nil, tree)

	return navContainer
}
