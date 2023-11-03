package ui_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/billiem/seren-management/src/ui"
)

func TestGetListableURI(t *testing.T) {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	baseTestDataDir := filepath.Join(workingDir, "../../test_data")

	var d = &ui.Data{}
	d.Config = &helpers.Config{BaseDir: baseTestDataDir}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "existing dir",
			path:     filepath.Join(baseTestDataDir, "files_test_dir"),
			expected: filepath.Join(baseTestDataDir, "files_test_dir"),
		},
		{
			name:     "existing file",
			path:     filepath.Join(baseTestDataDir, "files_test_dir/files_test_file.txt"),
			expected: filepath.Join(baseTestDataDir, "files_test_dir"),
		},
		{
			name:     "non-existing file, existing parent dir",
			path:     filepath.Join(baseTestDataDir, "files_test_dir/files_test_file_non_existent.txt"),
			expected: filepath.Join(baseTestDataDir, "files_test_dir"),
		},
		{
			name:     "non-existing file, non-existing parent dir",
			path:     filepath.Join(baseTestDataDir, "files_test_dir_non_existent/files_test_file_non_existent.txt"),
			expected: filepath.Join(baseTestDataDir),
		},
	}

	d.Config = &helpers.Config{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := d.GetListableURI(tt.path)

			fmt.Println(uri)
			// if uri.Path() != tt.expected {
			// 	t.Errorf("getListableURI(%s) returned %v, expected %v", tt.path, uri, tt.expected)
			// }
		})
	}
}
