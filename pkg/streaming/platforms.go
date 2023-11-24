package streaming

import "github.com/billiem/seren-management/pkg/helpers"

type StreamingPlatform interface {
	GetPlaylist() error
}

type GetPlaylistOpts interface {
	Build(helpers.Config) StreamingPlatform
}
