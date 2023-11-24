package collection

import (
	"strings"

	"github.com/billiem/seren-management/pkg/helpers"
)

type CollectionPlatform interface {
	ReadCollection() error
	UpdateCollection() error
}

type ReadCollectionOpts interface {
	Build(helpers.Config) CollectionPlatform
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
