package collection

import (
	"fmt"

	"github.com/billiem/seren-management/pkg/helpers"
)

/*
Contains a selection of utilities for managing a Traktor collection
*/

type ReadTraktorOpts struct {
	CollectionPath string
}

func (o ReadTraktorOpts) Build(cfg helpers.Config) PlatformCollection {
	collectionPath := cfg.TraktorCollectionPath

	if o.CollectionPath != "" {
		collectionPath = o.CollectionPath
	}

	return Traktor{
		CollectionPath: collectionPath,
	}
}

type Traktor struct {
	CollectionPath string
}

func (c Traktor) String() string {
	return "Traktor"
}

func (t Traktor) ReadCollection() {
	fmt.Println("read traktor collection")
}
