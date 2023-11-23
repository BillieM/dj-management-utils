package collection

import (
	"encoding/xml"
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

	coll := &NML{}

	err = xml.Unmarshal(data, coll)

	if err != nil {
		return err
	}

	for _, track := range coll.COLLECTION.ENTRY {
		fmt.Println(track)
	}

	return nil
}

func (t Traktor) WriteCollection() error {
	fmt.Println("write traktor collection")

	return nil
}
