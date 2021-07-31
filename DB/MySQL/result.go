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

type result struct {
	res		db.Result
	lastInsertId	*int64
	rowsAffected	*int64
}

// Make a new one of these!
func NewResult(res db.Result) ResultIfc {
	r := result{
		res:	res,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// ResultIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *result) GetLastInsertId() (*int64, error) {
	// If we don't have this cached yet, pull it from the result
	if nil == (*r).res {
		v, err := (*r).res.LastInsertId()
		if nil != err { return nil, err }
		(*r).lastInsertId = &v
	}
	return (*r).lastInsertId, nil
}

func (r *result) GetRowsAffected() (*int64, error) {
	// If we don't have this cached yet, pull it from the result
	if nil == (*r).res {
		v, err := (*r).res.RowsAffected()
		if nil != err { return nil, err }
		(*r).rowsAffected = &v
	}
	return (*r).rowsAffected, nil
}