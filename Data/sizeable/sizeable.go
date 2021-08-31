package sizeable

import (
	"bytes"
	"encoding/gob"
)

/*
Non-trivial data structures, such as an interface{} or a struct with nested objects, have no direct means of determining
their total data size. There are many possible reasons that one might need to know the size of such a structure. While
it is possible to use a new gob.NewEncoder().Encode() to Marshall the structure into a byte stream and then measure the
result, this would take time to process since the data must move/transform to discover the result. For applications that
have control over their own such structures, we provide this standardized interface for the bearer of such a structure to
determine its size such that the structure itself which implements this interface would contain the logic necessary to
efficiently determine its own data size without such a heavy marshaling operation being needed.

TODO:
 * Use the gob.Encode() recursive crawler as a baseline to traverse the value's structure and tally element sizes as
   opposed to transferring data: https://cs.opensource.google/go/go/+/refs/tags/go1.17:src/encoding/gob/encode.go
*/

type Sizeable interface {
	Size() int64
}

func Size(value interface{}) int64 {
	// If the value is Sizeable...
	if sizeableValue, ok := value.(Sizeable); ok { return sizeableValue.Size() }

	// If we can gob Encode it...
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(value); err == nil { return b.Len() }

	// Otherwise the size is not knowable
	return -1
}
