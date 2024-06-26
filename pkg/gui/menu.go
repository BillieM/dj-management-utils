package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
makeNavMenu builds the navigation menu on the left side of the application

the navigation menu is a tree object and is dynamically built from the operations list inside data.go
*/

func (e *guiEnv) makeNavMenu(contentStack *fyne.Container) *widget.Tree {

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return e.viewIndices[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := e.viewIndices[uid]
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("ExampleNodeTemplateText")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			op, ok := e.views[uid]
			if !ok {
				e.logger.FatalError(helpers.ErrOperationNotFound)
				return
			}
			node.(*widget.Label).SetText(op.name)
		},
		OnSelected: func(uid string) {
			op, ok := e.views[uid]
			if !ok {
				e.logger.FatalError(helpers.ErrOperationNotFound)
				return
			}
			if e.isBusy() {
				return
			}
			// if changing view, clear the external writer
			e.termSink.SetDiscard()
			e.setMainContent(contentStack, op)
		},
	}

	tree.OpenAllBranches()

	return tree
}
