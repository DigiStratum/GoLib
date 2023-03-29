package dependencies

/*

A primitive DependencyInstance Collection

This may seem a little cheeseball, but the idea here is that we can pack up a bunch of dependencies
into this collection, and then pass the entire collection across several layers until it reaches
the intended destination, and finally unpack the collection in the same way that we would inject
dependencies.

*/

type DependencyInstanceCollectionIfc interface {
	GetDependencyInstances() []DependencyInstanceIfc
}

type dependencyInstanceCollection struct {
	depinstcol		[]DependencyInstanceIfc
}

// --------------------------------------------------------------------------------------------------
// Factory Functions
// --------------------------------------------------------------------------------------------------

func NewDependencyInstanceCollection(depinst ...DependencyInstanceIfc) *dependencyInstanceCollection {
	dic := dependencyInstanceCollection{
		depinstcol:	make([]DependencyInstanceIfc, 0),
	}

	for _, di := range depinst { dic.depinstcol = append(dic.depinstcol, di) }

	return &dic
}

// --------------------------------------------------------------------------------------------------
// DependencyInstanceCollectionIfc
// --------------------------------------------------------------------------------------------------

// Return a slice of the DependencyInstances so that the recipient can just depinst... unpack normally
func (r *dependencyInstanceCollection) GetDependencyInstances() []DependencyInstanceIfc {
	return r.depinstcol
}

