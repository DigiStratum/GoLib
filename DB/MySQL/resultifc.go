package mysql

/*
Result Interface - When a query is created, it needs a prototype (ResultIfc) struct with Property Pointers to Scan() result data into
*/

type PropertyPointers []interface{}

// Result public interface
type ResultIfc interface {
	ZeroClone() (ResultIfc, PropertyPointers)
}

type ResultSet []ResultIfc

