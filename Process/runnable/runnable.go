package runnable

/*
A standardized interface for runnable background worker threads
*/

type RunnableIfc interface {
	Run()
	IsRunning() bool
	Stop()
}