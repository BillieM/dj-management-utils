package operations

import (
	"testing"

	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
	"github.com/google/go-cmp/cmp"
)

func TestBuildConvertTrack(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		outDirPath  string
		expected    ConvertTrack
		expectedErr error
	}{
		{
			name: "valid path",
			path: "/path/to/file.wav",
			expected: ConvertTrack{
				ID:   0,
				Name: "file",
				OriginalFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".wav",
						FullPath:      "/path/to/file.wav",
					},
				},
				NewFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/file.mp3",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "valid path with outDirPath",
			path:       "/path/to/file.wav",
			outDirPath: "/path/to/output/",
			expected: ConvertTrack{
				ID:   1,
				Name: "file",
				OriginalFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".wav",
						FullPath:      "/path/to/file.wav",
					},
				},
				NewFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/output/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/output/file.mp3",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:        "invalid path (no extension)",
			path:        "/path/to/nonexistent/file",
			expected:    ConvertTrack{},
			expectedErr: helpers.ErrMissingRequiredFields,
		},
		{
			name:        "invalid path (no file name/extension)",
			path:        "/path/to/nonexistent/",
			expected:    ConvertTrack{},
			expectedErr: helpers.ErrMissingRequiredFields,
		},
		{
			name:        "invalid path (no dir path)",
			path:        "file.wav",
			expected:    ConvertTrack{},
			expectedErr: helpers.ErrMissingRequiredFields,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := buildConvertTrack(i, tt.path, tt.outDirPath)

			if !helpers.ErrorContains(actualErr, tt.expectedErr) {
				t.Errorf("expected %v, but got %v", tt.expectedErr, actualErr)
			}

			// compare actual and expected ConvertTrack structs here
			if !cmp.Equal(actual, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
