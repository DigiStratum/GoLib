package chrono

/*
Abstraction of time-based considerations for basic timeStamp handling so that we can more easily
refactor later without bothering consumers. All times will be based on UTC, unix timeStamps. Initial
implementation uses Go runtime environment which could vary from one host to the next.

Defaults to local TimeSource
*/

import (
	"math"
	"time"
)

type TimeStampIfc interface {
	Add(offset int64) *TimeStamp
	Compare(ts TimeStamp) int
	Diff(ts TimeStamp) int64
}

type TimeStamp struct {
	timeStamp	int64
	timeSource	TimeSourceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimeStamp() *TimeStamp {
	timeStamp := TimeStamp{}
	return timeStamp.SetTimeSource(NewTimeSource())
}

// -------------------------------------------------------------------------------------------------
// TimeStampIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Chainable
func (r *TimeStamp) SetTimeSource(timeSource NewTimeSourceIfc) *TimeStamp {
	r.timeSource = timeSource
	r.timeStamp = timeSource.Now()
	return r
}

// Chainable
func (r *TimeStamp) Add(offset int64) *TimeStamp {
	r.timeStamp += offset
	return r
}

func (r TimeStamp) Compare(ts TimeStamp) int {
	if r.timeStamp == ts.timeStamp { return 0 }
	if r.timeStamp < ts.timeStamp { return -1 }
	return 1
}

func (r TimeStamp) Diff(ts TimeStamp) int64 {
	return r.timeStamp - ts.timeStamp
}
