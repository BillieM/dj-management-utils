package internal_test

import (
	"testing"

	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
	"github.com/google/go-cmp/cmp"
)

func TestSplitDirPath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     internal.FileInfo
		err      error
	}{
		{
			name:     "file with extension",
			filePath: "H:/Music/processed/funky cool song.mp3",
			want: internal.FileInfo{
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
			want: internal.FileInfo{
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
			want: internal.FileInfo{
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
			want: internal.FileInfo{
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
			want: internal.FileInfo{
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
			want: internal.FileInfo{
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
			want:     internal.FileInfo{},
			err:      helpers.ErrNoMatchesFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo, err := internal.SplitFilePath(tt.filePath)

			if !helpers.ErrorContains(err, tt.err) {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}

			if !cmp.Equal(fileInfo, tt.want) {
				t.Errorf("Expected %s, got %s", tt.want, fileInfo)
			}
		})
	}
}

func TestBuildFullPath(t *testing.T) {
	tests := []struct {
		name     string
		fileInfo internal.FileInfo
		want     string
	}{
		{
			name: "trailing slash on dir",
			fileInfo: internal.FileInfo{
				DirPath:       "H:/Music/processed/",
				FileName:      "funky cool song",
				FileExtension: ".mp3",
			},
			want: "H:/Music/processed/funky cool song.mp3",
		},
		{
			name: "no trailing slash on dir",
			fileInfo: internal.FileInfo{
				DirPath:       "/mnt/h/Music/processed",
				FileName:      "funky cool song",
				FileExtension: ".wav",
			},
			want: "/mnt/h/Music/processed/funky cool song.wav",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := tt.fileInfo.BuildFullPath()

			if fullPath != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, fullPath)
			}
		})
	}
}
