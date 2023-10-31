package main

import (
	"github.com/billiem/seren-management/src/helpers"
	"github.com/billiem/seren-management/src/ui"
)

func main() {

	c, err := helpers.LoadConfig()

	if err != nil {
		helpers.HandleFatalError(err)
	}

	ui.Entry(c)
}
