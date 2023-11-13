package operations

import "github.com/billiem/seren-management/src/helpers"

/*
getStemPaths gets all of the files in the provided directory which should be converted to stems based on the config

if recursion is true, will also get files in subdirectories
*/
func getStemPaths(cfg helpers.Config, inDirPath string, recursion bool) ([]string, error) {
	stemPaths, err := helpers.GetFilesInDir(inDirPath, recursion)
	if err != nil {
		return nil, err
	}
	var validStemPaths []string
	for _, path := range stemPaths {
		if helpers.IsExtensionInArray(path, cfg.ExtensionsToConvertToStems) {
			validStemPaths = append(validStemPaths, path)
		}
	}
	return validStemPaths, nil
}
