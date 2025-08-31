package oauth2

/*

Extend the Timestamp library with AccessToken specific concerns. This avoids adding oauth2-specific
implementation to the more generalized Timestamp, but still leverages Timestamps as the base
implementation.

All times will be based on UTC, unix timestamps. Initial implementation uses Go runtime environment
which could vary from one host to the next. A centralized time source, such as MySQL database, may
end up being more reliable (?)

TODO:
 * Provide method to replace the TimeSource with a custom one as needed

*/

import (
	chrono "github.com/DigiStratum/GoLib/Chrono"
)

type TokenTimeIfc interface {
	Now() *chrono.TimeStamp
	ExpiresAt(expiresIn int64) *chrono.TimeStamp
	ExpiresIn(expiresAt *chrono.TimeStamp) int64
	IsExpired(expiresAt *chrono.TimeStamp) bool
}

type tokenTime struct {
	timeSource *chrono.TimeSource
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTokenTime() *tokenTime {
	return &tokenTime{
		timeSource: chrono.NewTimeSource(),
	}
}

// -------------------------------------------------------------------------------------------------
// TokenTimeIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get the current time
func (r *tokenTime) Now() *chrono.TimeStamp {
	return r.timeSource.Now()
}

// Get the time that this token expires at as now + expiresIn seconds
func (r *tokenTime) ExpiresAt(expiresIn int64) *chrono.TimeStamp {
	return r.Now().Add(expiresIn)
}

// Get the seconds remaining before expiresAt expires; 0 if expired!
func (r *tokenTime) ExpiresIn(expiresAt *chrono.TimeStamp) int64 {
	t := expiresAt.Diff(r.Now())
	if 0 < t {
		return t
	}
	return 0
}

func (r *tokenTime) IsExpired(expiresAt *chrono.TimeStamp) bool {
	return expiresAt.Diff(r.Now()) == 0
}
