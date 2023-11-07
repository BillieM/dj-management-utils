package operations

import (
	"testing"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/google/go-cmp/cmp"
)

func TestBuildConvertTrack(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		outDirPath  string
		expected    ConvertTrack
		expectedErr string
	}{
		{
			name: "valid path",
			path: "/path/to/file.wav",
			expected: ConvertTrack{
				Track: Track{
					Name: "file",
				},
				OriginalFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".wav",
						FullPath:      "/path/to/file.wav",
					},
				},
				NewFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/file.mp3",
					},
				},
			},
			expectedErr: "",
		},
		{
			name:       "valid path with outDirPath",
			path:       "/path/to/file.wav",
			outDirPath: "/path/to/output/",
			expected: ConvertTrack{
				Track: Track{
					Name: "file",
				},
				OriginalFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/",
						FileName:      "file",
						FileExtension: ".wav",
						FullPath:      "/path/to/file.wav",
					},
				},
				NewFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/output/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/output/file.mp3",
					},
				},
			},
			expectedErr: "",
		},
		{
			name:        "invalid path (no extension)",
			path:        "/path/to/nonexistent/file",
			expected:    ConvertTrack{},
			expectedErr: "missing required fields",
		},
		{
			name:        "invalid path (no file name/extension)",
			path:        "/path/to/nonexistent/",
			expected:    ConvertTrack{},
			expectedErr: "missing required fields",
		},
		{
			name:        "invalid path (no dir path)",
			path:        "file.wav",
			expected:    ConvertTrack{},
			expectedErr: "missing required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, actualErr := buildConvertTrack(tt.path, tt.outDirPath)

			if actualErr != nil && actualErr.Error() != tt.expectedErr {
				t.Errorf("expected error %v, but got %v", tt.expectedErr, actualErr)
			}

			// compare actual and expected ConvertTrack structs here
			if !cmp.Equal(actual, tt.expected) {
				t.Errorf("expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
