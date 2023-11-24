package streaming

import "github.com/billiem/seren-management/pkg/helpers"

type StreamingPlatform interface {
	GetPlaylist() error
}

type GetPlaylistOpts interface {
	Build(helpers.Config) StreamingPlatform
}

type Playlist interface {
	NumTracks() int
}

type Track interface {
	Name() string
	SetName(string)
}
