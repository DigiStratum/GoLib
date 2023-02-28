package runnable

/*
A standardized interface for runnable background worker threads

FIXME: This overlaps with StartableIfc, and potentially StoppableIfc; Should these be one and the
same? Which paradigm do we actually want?
*/

type RunnableIfc interface {
	Run()
	IsRunning() bool
	Stop()
}
