package gui

import (
	"sync"
)

func (e *guiEnv) loadApp(appLoaded func()) {

	opEnv := e.opEnv()

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		opEnv.CheckLocalPaths()
	}()

	go func() {
		defer wg.Done()
		opEnv.IndexCollections()
	}()

	go func() {
		defer wg.Done()
		opEnv.IndexLocalFolders()
	}()

	wg.Wait()

	appLoaded()
}
