package mysql

/*
A SQL Query is Runnable if it implements the SQLQueryIfc interface. SQLQuery is a default implementation
of this that supports a raw query pass-through. But this paves the way to implement more advanced ORM
style query builders independent of the Query runner.

TODO:
 * Add support for literal substitutions here that could get us things like variable table names and
   other unescaped values for query variance
 * query '???' placeholders; because args is interface{}, we can pass whatever we want for this.
   let's declare that any arg which is an array will be treated as a set to expand. we will expect
   that there be the same number of sets as there are '???' placeholders. We will replace the  Nth
   '???' placeholder in the order that they appear with the count of '?,?,...' placeholders that
   matches the Len() of the Nth array. If the count of arrays supplied does not match the count of
   '???' placeholders, then error

*/

import (
	"fmt"
)

type SQLQueryIfc interface {
	// Resolve the query as a string or else a non-nil error
	Resolve(args ... interface{}) (string, error)
}

type SQLQuery struct {
	query		string
}

func NewSQLQuery(query string) *SQLQuery {
	return &SQLQuery{
		query:		query,
	}
}

// -------------------------------------------------------------------------------------------------
// SQLQueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

/*
Resolve this query object as a flat SQL statement ready for execution
Supplied args are same as those that will be fed to a prepared statement. These can be used to hint
our resolver here with macro tokens that will sub in the correct number of placeholders for prepared
statements, etc.

Using Pointer Receiver because we we don't want the caller to have to dereference a pointer produced
by NewSQLQuery().
*/
func (r *SQLQuery) Resolve(args ... interface{}) (string,  error) {
	if nil == r { return "", fmt.Errorf("SQLQuery is nil, nothing to resolve") }
	if len(r.query) == 0 { return "", fmt.Errorf("Query is empty") }
	// TODO: do some basic syntax/token/placeholder checks on query
	sql, err := r.resolveMacroTokens(args)
	if nil != err { return "", err }
	return *sql, nil
}

// -------------------------------------------------------------------------------------------------
// SQLQuery Private Implementation
// -------------------------------------------------------------------------------------------------

// Placeholder to support resolving magic expander tags, etc within our query
func (r SQLQuery) resolveMacroTokens(args ... interface{}) (*string, error) {
	protoQuery := r.query
	// TODO: Define and implement some useful macro token(s) such as '???' to replace a
	// some count of args value(s) with a matching count of '?' placeholders, etc.
        finalQuery := protoQuery
        return &finalQuery, nil
}

