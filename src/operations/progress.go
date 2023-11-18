package operations

type progress struct {
	completedTracks int
	totalTracks     int

	numSteps int

	inProcess map[int]int
}

func (p *progress) step(id int) float64 {
	p.inProcess[id]++
	return p.progress()
}

func (p *progress) complete(id int) float64 {
	delete(p.inProcess, id)
	p.completedTracks++
	return p.progress()
}

func (p *progress) progress() float64 {
	totalProgress := 0.0
	totalProgress += float64(p.completedTracks) / float64(p.totalTracks)
	for _, v := range p.inProcess {
		totalProgress += float64(v/p.numSteps) / float64(p.totalTracks)
	}
	return totalProgress
}

func buildProgress(totalTracks int, numSteps int) progress {
	return progress{
		completedTracks: 0,
		totalTracks:     totalTracks,
		numSteps:        numSteps,
		inProcess:       make(map[int]int),
	}
}
