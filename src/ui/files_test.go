package ui_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/billiem/seren-management/src/ui"
)

func TestGetClosestDir(t *testing.T) {

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
		baseDir  string
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
			path:     filepath.Join(baseTestDataDir, "files_test_dir/fake_file.txt"),
			expected: filepath.Join(baseTestDataDir, "files_test_dir"),
		},
		{
			name:     "non-existing file, non-existing parent dir",
			path:     filepath.Join(baseTestDataDir, "fake_dir/fake_file.txt"),
			expected: filepath.Join(baseTestDataDir),
		},
		{
			name:     "Many levels deep, returns BaseDir",
			path:     filepath.Join("/fake_dir_1/fake_dir2/fake_dir_3/fake_dir_4/file.txt"),
			expected: filepath.Join(baseTestDataDir),
			baseDir:  baseTestDataDir,
		},
		{
			name:     "Many levels deep, fake BaseDir, returns default /",
			path:     filepath.Join(baseTestDataDir, "fake_dir_1/fake_dir2/fake_dir_3/fake_dir_4/file.txt"),
			expected: filepath.Join("/"),
			baseDir:  filepath.Join("/fake_dir/weeoeoeoeo/"),
		},
	}

	d.Config = &helpers.Config{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d.Config.BaseDir = tt.baseDir
			var recursionCount int
			path := d.GetClosestDir(tt.path, &recursionCount)

			if path != tt.expected {
				t.Errorf("getListableURI(%s) returned %v, expected %v", tt.path, path, tt.expected)
			}
		})
	}
}
