package worker

/*
A standardized interface for background worker threads
*/

type WorkerIfc interface {
	func Run()
	func IsRunning()
	func Loop()
	func Stop()
}