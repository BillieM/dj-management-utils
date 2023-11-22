package collection

import (
	"fmt"
	"os"

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

func (t Traktor) ReadCollection() error {
	fmt.Println("read traktor collection", t.CollectionPath)

	// read xml
	data, err := os.ReadFile(t.CollectionPath)

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func (t Traktor) WriteCollection() error {
	fmt.Println("write traktor collection")

	return nil
}

type TraktorCollectionXML struct {
}
