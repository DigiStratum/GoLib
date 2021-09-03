package chrono

/*
A an abstract source of time that defaults to local system time

TODO:
 * A centralized time source, such as MySQL database, may end up being more reliable for cross-host
   comparisons (?)

*/

type TimeSourceIfc interface {
	Now() int64
	NewTimeStamp() *TimeStamp
}

type TimeSource struct {
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimeSource() *TimeSource {
	return &TimeSource{}
}

// -------------------------------------------------------------------------------------------------
// TimeSourceIfc Public Implementation
// -------------------------------------------------------------------------------------------------

func (r TimeSource) Now() int64 {
	return time.Now().Unix()
}

func (r TimeSource) NewTimeStamp() *TimeStamp {
	return NewTimeStamp(r)
}
