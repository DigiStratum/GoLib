package mysql

/*
Simple abstraction for database query results that prevent the caller from having to pull in Mysql database driver, etc.
This type of abstraction is important to support other types of database in the future.
*/

import (
	db "database/sql"
)

type ResultIfc interface {
	GetLastInsertId() (*int64, error)
	GetRowsAffected() (*int64, error)
}

type Result struct {
	res		db.Result
	lastInsertId	*int64
	rowsAffected	*int64
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewResult(res db.Result) *Result {
	if nil == res { return nil }
	return &Result{
		res:	res,
	}
}

// -------------------------------------------------------------------------------------------------
// ResultIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get the last insert id from our query result, if available
func (r *Result) GetLastInsertId() (*int64, error) {
	// If we don't have this cached yet, pull it from the result
	if nil == r.lastInsertId {
		v, err := r.res.LastInsertId()
		if nil != err { return nil, err }
		r.lastInsertId = &v
	}
	return r.lastInsertId, nil
}

// Get the affected row count from our query result, if available
func (r *Result) GetRowsAffected() (*int64, error) {
	// If we don't have this cached yet, pull it from the result
	if nil == r.rowsAffected {
		v, err := r.res.RowsAffected()
		if nil != err { return nil, err }
		r.rowsAffected = &v
	}
	return r.rowsAffected, nil
}
