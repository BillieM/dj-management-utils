package operations

/*
StemSeparationType is used to determine the type of stem output
*/
type StemSeparationType int

const (
	NotSelected StemSeparationType = iota // no value selected - needed for validation
	FourTrack                             // 4 .wav files for drums, bass, other, vocals
	Traktor                               // Traktor stems .stem.m4a
)

func (s StemSeparationType) String() string {
	return [...]string{"Not Selected", "4 Track", "Traktor"}[s]
}

func (s StemSeparationType) Check() bool {

	if s != FourTrack && s != Traktor {
		return false
	}
	return true
}
