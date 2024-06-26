package operations_test

import (
	"testing"

	"github.com/billiem/seren-management/pkg/operations/internal"
	stems "github.com/billiem/seren-management/pkg/operations/stems"
	"github.com/google/go-cmp/cmp"
)

func TestBuildStemTrack(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		outDirPath     string
		stemType       stems.StemSeparationType
		expectedOutput stems.StemTrack
		expectedError  error
	}{
		{
			name:       "Valid traktor path extraction with no outDirPath",
			path:       "/path/to/valid/file.mp3",
			outDirPath: "",
			stemType:   stems.Traktor,
			expectedOutput: stems.StemTrack{
				ID:         0,
				Name:       "file",
				StemDir:    "/path/to/valid/file/",
				SkipDemucs: false,
				StemsOnly:  false,
				OriginalFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "file",
						FileExtension: ".mp3",
						FullPath:      "/path/to/valid/file.mp3",
					},
				},
				OutFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "file",
						FileExtension: ".stem.m4a",
						FullPath:      "/path/to/valid/file.stem.m4a",
					},
				},
				BassFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "bass",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/bass.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				DrumsFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "drums",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/drums.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				OtherFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/path/to/valid/file/",
							FileName:      "other",
							FileExtension: ".mp3",
							FullPath:      "/path/to/valid/file/other.mp3",
						},
						DeleteOnFinish: true,
					},
				},
				VocalsFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
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
			stemType:   stems.FourTrack,
			expectedOutput: stems.StemTrack{
				ID:         1,
				Name:       "chicken",
				StemDir:    "/out/dir/path/chicken/",
				SkipDemucs: false,
				StemsOnly:  true,
				OriginalFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
						DirPath:       "/path/to/valid/",
						FileName:      "chicken",
						FileExtension: ".wav",
						FullPath:      "/path/to/valid/chicken.wav",
					},
				},
				OutFile: internal.AudioFile{},
				BassFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "bass",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/bass.wav",
						},
						DeleteOnFinish: false,
					},
				},
				DrumsFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "drums",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/drums.wav",
						},
						DeleteOnFinish: false,
					},
				},
				OtherFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
							DirPath:       "/out/dir/path/chicken/",
							FileName:      "other",
							FileExtension: ".wav",
							FullPath:      "/out/dir/path/chicken/other.wav",
						},
						DeleteOnFinish: false,
					},
				},
				VocalsFile: stems.StemFile{
					internal.AudioFile{
						FileInfo: internal.FileInfo{
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
			output, err := stems.BuildStemTrack(i, tt.path, tt.outDirPath, tt.stemType)

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
		want           stems.StemFile
	}{
		{
			name:           "bass stem delete on finish",
			baseStemDir:    "/path/to/stems/",
			fileName:       "bass",
			extension:      ".wav",
			deleteOnFinish: true,
			want: stems.StemFile{
				AudioFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
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
			want: stems.StemFile{
				AudioFile: internal.AudioFile{
					FileInfo: internal.FileInfo{
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
			got := stems.BuildStemFile(tt.baseStemDir, tt.fileName, tt.extension, tt.deleteOnFinish)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildStemFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
