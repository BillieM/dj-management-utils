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
	CollectionInPath  string
	CollectionOutPath string
}

func (o ReadTraktorOpts) Build(cfg helpers.Config) CollectionPlatform {
	var collectionInPath, collectionOutPath string

	if o.CollectionInPath == "" {
		collectionInPath = cfg.TraktorCollectionPath
	} else {
		collectionInPath = o.CollectionInPath
	}

	if o.CollectionOutPath == "" {
		collectionOutPath = fmt.Sprintf("%s_new.nml", helpers.RemoveFileExtension(collectionInPath))
	} else {
		collectionOutPath = o.CollectionOutPath
	}

	return &Traktor{
		CollectionInPath:  collectionInPath,
		CollectionOutPath: collectionOutPath,
		NML:               *new(NML),
	}
}

type Traktor struct {
	CollectionInPath  string
	CollectionOutPath string
	NML               NML
}

func (c Traktor) String() string {
	return "Traktor"
}

func (t Traktor) ReadCollection() error {
	err := t.loadCollection()

	if err != nil {
		return err
	}

	err = t.writeCollection()

	if err != nil {
		return err
	}

	return nil
}

func (t Traktor) UpdateCollection() error {
	fmt.Println("write traktor collection")

	return nil
}

func (t *Traktor) loadCollection() error {

	fmt.Println("load collection", t.CollectionInPath)

	// read xml
	data, err := os.ReadFile(t.CollectionInPath)

	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, &t.NML)

	if err != nil {
		return err
	}

	return nil
}

func (t *Traktor) writeCollection() error {

	fmt.Println("write collection", t.CollectionOutPath)

	collData, err := xml.MarshalIndent(t.NML, "", "  ")

	if err != nil {
		return err
	}

	writeData := []byte(traktorXMLHeader() + string(collData))

	err = os.WriteFile(t.CollectionOutPath, writeData, 0644)

	if err != nil {
		return err
	}

	return nil
}

func traktorXMLHeader() string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n"
}
