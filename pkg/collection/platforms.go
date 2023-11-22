package collection

import (
	"strings"
)

type PlatformCollection interface {
	ReadCollection() error
	WriteCollection() error
}

var (
	validPlatforms = []string{
		"traktor",
		"rekordbox",
		"serato",
	}
)

func CommaSeparatedPlatforms() string {
	return strings.Join(validPlatforms, ", ")
}
