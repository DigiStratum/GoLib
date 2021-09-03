package chrono

/*
Abstraction of time-based considerations for basic timestamp handling so that we can more easily
refactor later without bothering consumers. All times will be based on UTC, unix timestamps. Initial
implementation uses Go runtime environment which could vary from one host to the next.

TODO:
 * A centralized time source, such as MySQL database, may end up being more reliable for cross-host
   comparisons (?)

*/

import (
	"math"
	"time"
)

type TimestampIfc interface {
	Add(offset int64) *Timestamp
	Compare(ts Timestamp) int
	Diff(ts Timestamp) int64
}

type Timestamp struct {
	timestamp	int64
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimestamp() *Timestamp {
	return &Timestamp{
		timestamp:	time.Now().Unix(),
	}
}

// -------------------------------------------------------------------------------------------------
// TimestampIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Timestamp) Add(offset int64) *Timestamp {
	r.timestamp += offset
	return r
}

func (r Timestamp) Compare(ts Timestamp) int {
	if r.timestamp == ts.timestamp { return 0 }
	if r.timestamp < ts.timestamp { return -1 }
	return 1
}

func (r Timestamp) Diff(ts Timestamp) int64 {
	return r.timestamp - ts.timestamp
}
