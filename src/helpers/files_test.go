package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

// TODO: Convert these tests to use generic paths
// TODO: Add tests for the other functions in files.go

func TestReplaceTrackExtension(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		newExtension string
		extensions   []string
		want         string
	}{
		{
			name:         "replace mp3 with stem.m4a",
			filePath:     "H:/Music/processed/funky cool song.mp3",
			newExtension: ".stem.m4a",
			extensions:   []string{"mp3"},
			want:         "H:/Music/processed/funky cool song.stem.m4a",
		},
		{
			name:         "replace aiff with mp3",
			filePath:     "H:/Music/processed/funky cool song.aiff",
			newExtension: ".mp3",
			extensions:   []string{"aiff"},
			want:         "H:/Music/processed/funky cool song.mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPath := ReplaceTrackExtension(tt.filePath, tt.newExtension, tt.extensions)

			if newPath != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, newPath)
			}
		})
	}
}

func TestIsExtensionInArray(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		extensions []string
		want       bool
	}{
		{
			name:       "extension in array",
			filename:   "cool song.mp3",
			extensions: []string{"mp3", "aiff", "wav"},
			want:       true,
		},
		{
			name:       "extension not in array",
			filename:   "cool song.m4a",
			extensions: []string{"mp3", "aiff", "wav"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExtensionInArray(tt.filename, tt.extensions); got != tt.want {
				t.Errorf("IsExtensionInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileExtensionFromFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
		err      string
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     ".mp3",
			err:      "",
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "",
			err:      "no file extension found",
		},
		{
			name:     "file without dir path",
			filePath: "funky cool song.mp3",
			want:     ".mp3",
			err:      "",
		},
		{
			name:     "file without dir path or extension",
			filePath: "funky cool song",
			want:     "",
			err:      "no file extension found",
		},
		{
			name:     "dir path only",
			filePath: "H:/Music/processed/",
			want:     "",
			err:      "no file extension found",
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     ".m4a",
			err:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileExtension, err := GetFileExtensionFromFilePath(tt.filePath)

			if err != nil && err.Error() != tt.err {
				t.Errorf("Unexpected error: %v", err)
			}

			if fileExtension != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, fileExtension)
			}
		})
	}
}

func TestGetDirPathFromFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
		err      string
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     "H:/Music/processed/",
			err:      "",
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "H:/Music/processed/",
			err:      "",
		},
		{
			name:     "directory path only (trailing /)",
			filePath: "H:/Music/processed/",
			want:     "H:/Music/processed/",
			err:      "",
		},
		{
			name:     "file name only",
			filePath: "funky cool song",
			want:     "",
			err:      "no directory path found",
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     "H:/tmp/testdir/",
			err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDirPathFromFilePath(tt.filePath)

			if err != nil && err.Error() != tt.err {
				t.Errorf("GetDirPathFromFilePath() error = %v, wantErr %v", err, tt.err)
			}

			if got != tt.want {
				t.Errorf("GetDirPathFromFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileNameFromFilePath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
		err      string
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     "funky cool song",
			err:      "",
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "funky cool song",
			err:      "",
		},
		{
			name:     "directory path only (trailing /)",
			filePath: "H:/Music/processed/",
			want:     "",
			err:      "no file name found",
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     "01 - funky cool song (feat. coolman)",
			err:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileNameFromFilePath(tt.filePath)

			if err != nil && err.Error() != tt.err {
				t.Errorf("GetFileNameFromFilePath() error = %v, wantErr %v", err, tt.err)
			}

			if got != tt.want {
				t.Errorf("GetFileNameFromFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClosestDir(t *testing.T) {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	baseTestDataDir := filepath.Join(workingDir, "../../test_data")

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var recursionCount int
			path, _ := GetClosestDir(tt.path, tt.baseDir, &recursionCount)

			if path != tt.expected {
				t.Errorf("getListableURI(%s) returned %v, expected %v", tt.path, path, tt.expected)
			}
		})
	}
}
