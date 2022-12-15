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
	Now() TimeStampIfc
	NowUnixTimeStamp() int64
	NowUnixTimeStampMilli() int64
}

type timeSource struct {
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimeSource() *timeSource {
	return &timeSource{}
}

// -------------------------------------------------------------------------------------------------
// TimeSourceIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r timeSource) Now() TimeStampIfc {
	return NewTimeStamp(r)
}

func (r timeSource) NowUnixTimeStamp() int64 {
	return time.Now().Unix()
}

func (r timeSource) NowUnixTimeStampMilli() int64 {
	return time.Now().UnixMilli()
}

