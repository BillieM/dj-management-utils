package iwidget

import "github.com/billiem/seren-management/pkg/helpers"

/*
Contains a base struct for other internal widget structs,
this is primarily used to simplify logging
*/

type Base struct {
	Logger helpers.SerenLogger
}
