package operations

type progress struct {
	completed int
	total     int

	stepsPer int

	inProcess map[int]int
}

func (p *progress) step(id int) float64 {
	p.inProcess[id]++
	return p.value()
}

func (p *progress) complete(id int) float64 {
	delete(p.inProcess, id)
	p.completed++
	return p.value()
}

func (p *progress) value() float64 {
	totalProgress := 0.0
	totalProgress += float64(p.completed) / float64(p.total)
	for _, v := range p.inProcess {
		totalProgress += (float64(v) / float64(p.stepsPer) / float64(p.total))
	}
	return totalProgress
}

func buildProgress(total int, stepsPer int) progress {
	return progress{
		completed: 0,
		total:     total,
		stepsPer:  stepsPer,
		inProcess: make(map[int]int),
	}
}
