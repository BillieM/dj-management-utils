package internal

/*
progress is a struct used to track the progress of a process that may have multiple steps, for example, a stem separation process,
the value can then be used to e.g. update a progress bar
*/
type Progress struct {
	completed int
	total     int

	stepsPer int

	inProcess map[int]*int
}

func (p *Progress) Step(id int) float64 {
	*p.inProcess[id]++
	return p.value()
}

func (p *Progress) Complete(id int) float64 {
	delete(p.inProcess, id)
	p.completed++
	return p.value()
}

func (p *Progress) value() float64 {
	totalProgress := 0.0
	totalProgress += float64(p.completed) / float64(p.total)
	for _, v := range p.inProcess {
		totalProgress += (float64(*v) / float64(p.stepsPer) / float64(p.total))
	}
	return totalProgress
}

/*
buildProgress returns a progress struct with the given total number of processes and the number of steps per process

This is used to track the progress of a process that may have multiple steps, for example, a stem separation process,
the value can then be used to e.g. update a progress bar
*/
func BuildProgress(total int, stepsPer int) Progress {
	processMap := make(map[int]*int)

	for i := 0; i < total; i++ {
		processMap[i] = new(int)
	}

	return Progress{
		completed: 0,
		total:     total,
		stepsPer:  stepsPer,
		inProcess: processMap,
	}
}
