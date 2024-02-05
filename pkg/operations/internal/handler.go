package internal

/*
Provides functionality surrounding the operation 'handler'

This handler is responsible for communicating output from operations back to the
interface used to trigger the operation, be it the GUI or the CLI

Handlers are registered from the GUI/CLI prior to the operation being triggered, and
are then attached to the operation environment (OpEnv) which the operation is a method of

The operation will then be launched, at which point it will call the handler upon
failure, success, or just to provide progress updates
*/

/*
OperationHandler is a struct that is used to handle the output of an operation
back to the user interface
*/
type OperationHandler struct {
	progressTracker *progressTracker
	progressFunc    func(float64)
	successFunc     func(map[string]any)
	errorFunc       func(error)
}

func BuildOperationHandler(
	progressFunc func(float64),
	successFunc func(map[string]any),
	errorFunc func(error),
) *OperationHandler {
	return &OperationHandler{
		progressFunc: progressFunc,
		successFunc:  successFunc,
		errorFunc:    errorFunc,
	}
}

/*
BuildProgressTracker builds the underlying progress tracker for the operation,
this is used to generate float64 values representing the progress of the operation
*/
func (o *OperationHandler) BuildProgressTracker(totalProcs int, stepsPer int) {
	o.progressTracker = buildProgressTracker(totalProcs, stepsPer)
}

/*
ProcessStep increments the progress of the process with the given id by 1
and then calls the given progressFunc with the operations new total progress
*/
func (o *OperationHandler) ProcessStep(id int) {
	o.progressFunc(
		o.progressTracker.step(id),
	)
}

/*
ProcessComplete removes the process with the given id from the progress tracker
and then calls the given progressFunc with the operations new total progress
*/
func (o *OperationHandler) ProcessComplete(id int) {
	o.progressFunc(
		o.progressTracker.complete(id),
	)
}

func (o *OperationHandler) FinishError(err error) {
	o.progressFunc(1)
	o.errorFunc(err)
}

func (o *OperationHandler) FinishSuccess(data map[string]any) {
	o.progressFunc(1)
	o.successFunc(data)
}
