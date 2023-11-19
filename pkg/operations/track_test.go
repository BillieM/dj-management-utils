package operations

import (
	"testing"

	"github.com/billiem/seren-management/pkg/helpers"
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
				ID: 0,
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
			expectedErr: nil,
		},
		{
			name:       "valid path with outDirPath",
			path:       "/path/to/file.wav",
			outDirPath: "/path/to/output/",
			expected: ConvertTrack{
				ID: 1,
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

func TestBuildStemTrack(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		outDirPath     string
		stemType       StemSeparationType
		expectedOutput StemTrack
		expectedError  error
	}{
		{
			name:       "Valid traktor path extraction with no outDirPath",
			path:       "/path/to/valid/file.mp3",
			outDirPath: "",
			stemType:   Traktor,
			expectedOutput: StemTrack{
				ID: 0,
				Track: Track{
					Name: "file",
				},
				StemDir:    "/path/to/valid/file/",
				SkipDemucs: false,
				StemsOnly:  false,
				OriginalFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/valid/file.mp3",
					},
				},
				OutFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "file",
						FileExtension: ".stem.m4a",
						FullPath:      "/path/to/valid/file.stem.m4a",
					},
				},
				BassFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "bass",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/bass.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				DrumsFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "drums",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/drums.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				OtherFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "other",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/other.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				VocalsFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "vocals",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/vocals.mp3",
						},
						DeleteOnFinish: true,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:       "Valid FourTrack path extraction with outDirPath",
			path:       "/path/to/valid/chicken.wav",
			outDirPath: "/out/dir/path/",
			stemType:   FourTrack,
			expectedOutput: StemTrack{
				ID: 1,
				Track: Track{
					Name: "chicken",
				},
				StemDir:    "/out/dir/path/chicken/",
				SkipDemucs: false,
				StemsOnly:  true,
				OriginalFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "chicken",
						FileExtension: ".wav",
						FullPath:      "/path/to/valid/chicken.wav",
					},
				},
				OutFile: AudioFile{},
				BassFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "bass",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/bass.wav",
						},
						DeleteOnFinish: false,
					},
				},
				DrumsFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "drums",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/drums.wav",
						},
						DeleteOnFinish: false,
					},
				},
				OtherFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "other",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/other.wav",
						},
						DeleteOnFinish: false,
					},
				},
				VocalsFile: StemFile{
					AudioFile{
						FileInfo: helpers.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "vocals",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/vocals.wav",
						},
						DeleteOnFinish: false,
					},
				},
			},
			expectedError: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := buildStemTrack(i, tt.path, tt.outDirPath, tt.stemType)

			if diff := cmp.Diff(output, tt.expectedOutput); diff != "" {
				t.Errorf("buildStemTrack() output mismatch (-got +want):\n%s", diff)
			}

			if err != tt.expectedError {
				t.Errorf("buildStemTrack() error mismatch (got: %v, want: %v)", err, tt.expectedError)
			}
		})
	}
}

func TestBuildStemFile(t *testing.T) {
	tests := []struct {
		name           string
		baseStemDir    string
		fileName       string
		extension      string
		deleteOnFinish bool
		want           StemFile
	}{
		{
			name:           "bass stem delete on finish",
			baseStemDir:    "/path/to/stems/",
			fileName:       "bass",
			extension:      ".wav",
			deleteOnFinish: true,
			want: StemFile{
				AudioFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/stems/",
						FileName:      "bass",
						FileExtension: ".wav",
						FullPath:      "/path/to/stems/bass.wav",
					},
					DeleteOnFinish: true,
				},
			},
		},
		{
			name:           "drum stem don't delete on finish",
			baseStemDir:    "/path/to/stems/",
			fileName:       "drums",
			extension:      ".mp3",
			deleteOnFinish: false,
			want: StemFile{
				AudioFile: AudioFile{
					FileInfo: helpers.FileInfo{
						DirPath:       "/path/to/stems/",
						FileName:      "drums",
						FileExtension: ".mp3",
						FullPath:      "/path/to/stems/drums.mp3",
					},
					DeleteOnFinish: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildStemFile(tt.baseStemDir, tt.fileName, tt.extension, tt.deleteOnFinish)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildStemFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
