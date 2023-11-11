package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/google/go-cmp/cmp"
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
			newPath := helpers.ReplaceTrackExtension(tt.filePath, tt.newExtension, tt.extensions)

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
			if got := helpers.IsExtensionInArray(tt.filename, tt.extensions); got != tt.want {
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
		err      error
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     ".mp3",
			err:      nil,
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "",
			err:      helpers.ErrNoFileExtension,
		},
		{
			name:     "file without dir path",
			filePath: "funky cool song.mp3",
			want:     ".mp3",
			err:      nil,
		},
		{
			name:     "file without dir path or extension",
			filePath: "funky cool song",
			want:     "",
			err:      helpers.ErrNoFileExtension,
		},
		{
			name:     "dir path only",
			filePath: "H:/Music/processed/",
			want:     "",
			err:      helpers.ErrNoFileExtension,
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     ".m4a",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileExtension, err := helpers.GetFileExtensionFromFilePath(tt.filePath)

			if !helpers.ErrorContains(err, tt.err) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
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
		err      error
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     "H:/Music/processed/",
			err:      nil,
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "H:/Music/processed/",
			err:      nil,
		},
		{
			name:     "directory path only (trailing /)",
			filePath: "H:/Music/processed/",
			want:     "H:/Music/processed/",
			err:      nil,
		},
		{
			name:     "file name only",
			filePath: "funky cool song",
			want:     "",
			err:      helpers.ErrNoDirPath,
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     "H:/tmp/testdir/",
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.GetDirPathFromFilePath(tt.filePath)

			if !helpers.ErrorContains(err, tt.err) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
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
		err      error
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want:     "funky cool song",
			err:      nil,
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want:     "funky cool song",
			err:      nil,
		},
		{
			name:     "directory path only (trailing /)",
			filePath: "H:/Music/processed/",
			want:     "",
			err:      helpers.ErrNoFileName,
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want:     "01 - funky cool song (feat. coolman)",
			err:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.GetFileNameFromFilePath(tt.filePath)

			if !helpers.ErrorContains(err, tt.err) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
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
			path, _ := helpers.GetClosestDir(tt.path, tt.baseDir, &recursionCount)

			if path != tt.expected {
				t.Errorf("getListableURI(%s) returned %v, expected %v", tt.path, path, tt.expected)
			}
		})
	}
}
func TestSplitDirPath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     helpers.FileInfo
		err      error
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want: helpers.FileInfo{
				FullPath:      "H:/Music/processed/funky cool song.mp3",
				DirPath:       "H:/Music/processed/",
				FileName:      "funky cool song",
				FileExtension: ".mp3",
			},
			err: nil,
		},
		{
			name:     "file without extension",
			filePath: "H:/Music/processed/funky cool song",
			want: helpers.FileInfo{
				FullPath:      "H:/Music/processed/funky cool song",
				DirPath:       "H:/Music/processed/",
				FileName:      "funky cool song",
				FileExtension: "",
			},
			err: nil,
		},
		{
			name:     "file without dir path",
			filePath: "funky cool song.mp3",
			want: helpers.FileInfo{
				FullPath:      "funky cool song.mp3",
				DirPath:       "",
				FileName:      "funky cool song",
				FileExtension: ".mp3",
			},
			err: nil,
		},
		{
			name:     "file without dir path or extension",
			filePath: "funky cool song",
			want: helpers.FileInfo{
				FullPath:      "funky cool song",
				DirPath:       "",
				FileName:      "funky cool song",
				FileExtension: "",
			},
			err: nil,
		},
		{
			name:     "dir path only",
			filePath: "H:/Music/processed/",
			want: helpers.FileInfo{
				FullPath:      "H:/Music/processed/",
				DirPath:       "H:/Music/processed/",
				FileName:      "",
				FileExtension: "",
			},
			err: nil,
		},
		{
			name:     "file name with dot",
			filePath: "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
			want: helpers.FileInfo{
				FullPath:      "H:/tmp/testdir/01 - funky cool song (feat. coolman).m4a",
				DirPath:       "H:/tmp/testdir/",
				FileName:      "01 - funky cool song (feat. coolman)",
				FileExtension: ".m4a",
			},
			err: nil,
		},
		{
			name:     "empty string",
			filePath: "",
			want:     helpers.FileInfo{},
			err:      helpers.ErrNoMatchesFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo, err := helpers.SplitFilePath(tt.filePath)

			if !helpers.ErrorContains(err, tt.err) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}

			if !cmp.Equal(fileInfo, tt.want) {
				t.Errorf("Expected %s, got %s", tt.want, fileInfo)
			}
		})
	}
}
