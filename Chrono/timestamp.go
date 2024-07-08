package chrono

/*
Abstraction of time-based considerations for basic timeStamp handling with 1-second precisionso
that we can more easily refactor later without bothering consumers. All times will be based on UTC,
unix timeStamps. Initial implementation uses Go runtime environment which could vary from one host
to the next. Defaults to local TimeSource.

*/

//import ( "fmt" )

type TimeStampIfc interface {
	Add(offset int64) *timeStamp
	Compare(ts TimeStampIfc) int
	CompareToNow() int
	Diff(ts TimeStampIfc) int64
	DiffNow() int64

	IsForever() bool
	IsPast() bool
	IsFuture() bool

	ToUnixTimeStamp() int64
}

type timeStamp struct {
	timeStamp	int64
	timeSource	TimeSourceIfc
	isForever	bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTimeStamp(timeSource TimeSourceIfc) *timeStamp {
	if nil == timeSource { return nil }
	return NewFromUnixTimeStamp(timeSource, timeSource.NowUnixTimeStamp())
}

func NewFromUnixTimeStamp(timeSource TimeSourceIfc, unixTimeStamp int64) *timeStamp {
	// Require a TimeSource
	if nil == timeSource { return nil }
	ts := timeStamp{
		timeStamp: unixTimeStamp,
		timeSource: timeSource,
		isForever: false,
	}
	return &ts
}

func NewTimeStampForever() *timeStamp {
	ts := timeStamp{
		isForever:	true,
	}
	return &ts
}

// -------------------------------------------------------------------------------------------------
// TimeStampIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Chainable
func (r *timeStamp) Add(offset int64) *timeStamp {
	if ! r.isForever { r.timeStamp += offset }
	return r
}

func (r *timeStamp) Compare(ts TimeStampIfc) int {
	if r.isForever { return 1 }			// Forever is always in the future
	if t, ok := ts.(*timeStamp); ok {
		if r.timeStamp < t.timeStamp { return -1 }	// Past
		if r.timeStamp == t.timeStamp { return 0 }	// Present
	}
	return 1					// Future
}

func (r *timeStamp) CompareToNow() int {
	if r.isForever { return 1 }			// Forever is always in the future
	return r.Compare(r.timeSource.Now())
}

func (r *timeStamp) Diff(ts TimeStampIfc) int64 {
	if r.isForever { return 1 }			// Forever is always in the future
	if t, ok := ts.(*timeStamp); ok {
    //fmt.Printf("Diff = %d - %d = %d\n", r.timeStamp, ts.timeStamp, r.timeStamp - ts.timeStamp)
		return r.timeStamp - t.timeStamp
	}
	return 0
}

func (r *timeStamp) DiffNow() int64 {
	if r.isForever { return 1 }			// Forever is always in the future
	now := r.timeSource.Now()
	return r.Diff(now)
}

func (r *timeStamp) IsForever() bool {
	return r.isForever
}

func (r *timeStamp) IsPast() bool {
	if r.isForever { return false }			// Forever is never in the past
	diff := r.DiffNow()
	return (diff < 0)
}

func (r *timeStamp) IsFuture() bool {
	if r.isForever { return true }			// Forever is always in the future
	return (r.DiffNow() > 0)
}

func (r *timeStamp) ToUnixTimeStamp() int64 {
	if r.isForever { return 0 }			// Forever has no definite timestamp
	return r.timeStamp
}

