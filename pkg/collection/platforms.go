package collection

import (
	"strings"
)

type CollectionPlatform interface {
	ReadCollection() error
	UpdateCollection() error
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
