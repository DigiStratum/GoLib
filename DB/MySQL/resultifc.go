package mysql

type PropertyPointers []interface{}

type ResultIfc interface {
	ZeroClone() (ResultIfc, PropertyPointers)
}

type ResultSet []ResultIfc

