package chrono

/*
Abstraction of time-based considerations for basic timeStamp handling so that we can more easily
refactor later without bothering consumers. All times will be based on UTC, unix timeStamps. Initial
implementation uses Go runtime environment which could vary from one host to the next.

Defaults to local TimeSource
*/

type TimeStampIfc interface {
	Add(offset int64) *TimeStamp
	Compare(ts *TimeStamp) int
	CompareToNow() int
	Diff(ts *TimeStamp) int64
	DiffNow() int64

	IsForever() bool
	IsPast() bool
	IsFuture() bool

	ToUnixTimeStamp() int64
}

type TimeStamp struct {
	timeStamp	int64
	timeSource	TimeSourceIfc
	isForever	bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimeStamp(timeSource TimeSourceIfc) *TimeStamp {
	return NewFromUnixTimeStamp(timeSource, timeSource.NowUnixTimeStamp())
}

func NewFromUnixTimeStamp(timeSource TimeSourceIfc, unixTimeStamp int64) *TimeStamp {
	return &TimeStamp{
		timeStamp: unixTimeStamp,
		timeSource: timeSource,
	}
}

func NewTimeStampForever() *TimeStamp {
	return &TimeStamp{
		isForever:	true,
	}
}

// -------------------------------------------------------------------------------------------------
// TimeStampIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Chainable
func (r *TimeStamp) Add(offset int64) *TimeStamp {
	if ! r.isForever { r.timeStamp += offset }
	return r
}

func (r TimeStamp) Compare(ts *TimeStamp) int {
	if r.isForever { return 1 } // Forever is always in the future
	if r.timeStamp == ts.timeStamp { return 0 }
	if r.timeStamp < ts.timeStamp { return -1 }
	return 1
}

func (r TimeStamp) CompareToNow() int {
	if r.isForever { return 1 } // Forever is always in the future
	return r.Compare(r.timeSource.Now())
}

func (r TimeStamp) Diff(ts *TimeStamp) int64 {
	if r.isForever { return 1 } // Forever is always in the future
	return r.timeStamp - ts.timeStamp
}

func (r TimeStamp) DiffNow() int64 {
	if r.isForever { return 1 } // Forever is always in the future
	return r.Diff(r.timeSource.Now())
}

func (r TimeStamp) IsForever() bool {
	return r.isForever
}

func (r TimeStamp) IsPast() bool {
	return (r.DiffNow() == -1)
}

func (r TimeStamp) IsFuture() bool {
	return (r.DiffNow() == 1)
}

func (r TimeStamp) ToUnixTimeStamp() int64 {
	return r.timeStamp
}
