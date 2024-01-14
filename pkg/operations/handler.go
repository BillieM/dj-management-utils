package operations

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
OperationHandler is used to provide callbacks to the operations package
*/
type operationHandler struct {
	StepCallback     func(float64)
	FinishedCallback func(OperationFinishedInfo)
}

/*
OperationFinishedInfo provides a format for passing information about the completion of an operation
back to the interface that triggered it
*/
type OperationFinishedInfo struct {
	Data map[string]any // TODO: consider changing this to an interface
	Err  error
}

/*
RegisterOperationHandler registers an OperationHandler with the OpEnv
*/
func (e *OpEnv) RegisterOperationHandler(stepCallback func(float64), finishedCallback func(OperationFinishedInfo)) {
	e.operationHandler = operationHandler{
		StepCallback:     stepCallback,
		FinishedCallback: finishedCallback,
	}
}

/*
progress is called by an operation, it calls the StepCallback assigned to the OpEnv
*/
func (e *OpEnv) progress(stepInfo float64) {
	e.operationHandler.StepCallback(stepInfo)
}

/*
finish is called by an operation, it calls the FinishedCallback assigned to the OpEnv
*/
func (e *OpEnv) finish(finishedInfo OperationFinishedInfo) {
	e.operationHandler.FinishedCallback(finishedInfo)
}

/*
finishError is a helper function for generating an OperationFinishedInfo struct for an error
and then calling finish, triggering the FinishedCallback assigned to the OperationHandler
*/
func (e *OpEnv) finishError(err error) {
	e.finish(newOperationFinishedInfoError(err))
}

/*
finishSuccess is a helper function for generating an OperationFinishedInfo struct for success
and then calling finish, triggering the FinishedCallback assigned to the OperationHandler
*/
func (e *OpEnv) finishSuccess(data map[string]any) {
	e.finish(newOperationFinishedInfoSuccess(data))
}

/*
newFinishedError is a helper function for generating an OperationFinishedInfo struct
*/
func newOperationFinishedInfoError(err error) OperationFinishedInfo {
	return OperationFinishedInfo{
		Err: err,
	}
}

/*
 */
func newOperationFinishedInfoSuccess(data map[string]any) OperationFinishedInfo {
	return OperationFinishedInfo{
		Data: data,
		Err:  nil,
	}
}

/*
 */
// func newOperationProgressInfo(msg string) OperationProgressInfo {
// 	return OperationProgressInfo{
// 		Message:    msg,
// 		Importance: helpers.HighImportance,
// 	}
// }
