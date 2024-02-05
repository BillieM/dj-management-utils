package internal

/*
progress is a struct used to track the progress of a process that may have multiple steps, for example, a stem separation process,
the value can then be used to e.g. update a progress bar
*/
type progressTracker struct {
	completed int
	total     int

	stepsPer int

	// inProcess is a map of the id of the process to the current step of the process
	// we use a map of pointers to avoid concurrency issues
	inProcess map[int]*int
}

func (p *progressTracker) step(id int) float64 {
	*p.inProcess[id]++
	return p.value()
}

func (p *progressTracker) complete(id int) float64 {
	delete(p.inProcess, id)
	p.completed++
	return p.value()
}

func (p *progressTracker) value() float64 {
	totalProgress := 0.0
	totalProgress += float64(p.completed) / float64(p.total)
	for _, v := range p.inProcess {
		totalProgress += (float64(*v) / float64(p.stepsPer) / float64(p.total))
	}
	return totalProgress
}

/*
buildProgressTracker returns a progress struct with the given total number of processes and the number of steps per process

This is used to track the progress of a process that may have multiple steps, for example, a stem separation process,
the value can then be used to e.g. update a progress bar
*/
func buildProgressTracker(totalProcs int, stepsPer int) *progressTracker {
	processMap := make(map[int]*int)

	for i := 0; i < totalProcs; i++ {
		processMap[i] = new(int)
	}

	return &progressTracker{
		completed: 0,
		total:     totalProcs,
		stepsPer:  stepsPer,
		inProcess: processMap,
	}
}
