package ui

type OperationProcess struct {
	progressBar MyProgressBar
}

func (o OperationProcess) StepCallback(progress float64) {
	o.progressBar.updateProgressBar(progress)
}
