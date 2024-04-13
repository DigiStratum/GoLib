package fileio

type FileSetIfc interface {
}

type FileSet struct {
	files		[]*file
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewFileSet() *FileSet {
	r := FileSet{}
	return &r
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

func (r *FileSet) GetIterator() func () interface{} {
	/*
	// TODO: lifted from DB.ResultSet; update to iterate over this collection of r.files
        idx := 0
        var data_len = r.Len()
        return func () interface{} {
                // If we're done iterating, return do nothing
                if idx >= data_len { return nil }
                prev_idx := idx
                idx++
                return &r.results[prev_idx]
        }
	*/
	return nil
}
