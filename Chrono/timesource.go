package chrono

/*
A an abstract source of time that defaults to local system time

TODO:
 * A centralized time source, such as MySQL database, may end up being more reliable for cross-host
   comparisons (?)

*/

import (
	"time"
)

type TimeSourceIfc interface {
	Now() *TimeStamp
	NowUnixTimeStamp() int64
	NowUnixTimeStampMilli() int64
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
// TimeSourceIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r TimeSource) Now() *TimeStamp {
	return NewTimeStamp(r)
}

func (r TimeSource) NowUnixTimeStamp() int64 {
	return time.Now().Unix()
}

func (r TimeSource) NowUnixTimeStampMilli() int64 {
	return time.Now().UnixMilli()
}
