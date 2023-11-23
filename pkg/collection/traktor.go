package collection

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/k0kubun/pp"
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

	return &Traktor{
		CollectionPath: collectionPath,
	}
}

type Traktor struct {
	CollectionPath string
	*NML
}

func (c Traktor) String() string {
	return "Traktor"
}

func (t Traktor) ReadCollection() error {
	fmt.Println("read traktor collection", t.CollectionPath)
	err := t.loadTraktorCollection()

	if err != nil {
		return err
	}

	return nil

}

func (t Traktor) UpdateCollection() error {
	fmt.Println("write traktor collection")

	return nil
}

func (t *Traktor) loadTraktorCollection() error {

	// read xml
	data, err := os.ReadFile(t.CollectionPath)

	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, t.NML)

	if err != nil {
		return err
	}

	pp.Print(t.NML.PLAYLISTS.NODE)

	return nil
}

func (t *Traktor) writeCollection() {

}
