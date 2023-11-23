package streaming

type StreamingPlatform interface {
	GetPlaylists() error
}
