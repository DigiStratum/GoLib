package runnable

/*
A standardized interface for runnable background worker threads
*/

type RunnableIfc interface {
	func Run()
	func IsRunning()
	func Stop()
}