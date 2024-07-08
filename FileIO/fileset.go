package fileio

type FileSetIfc interface {
	AddFile(path string)
	Len() int
}

type fileSet struct {
	files		[]*file
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewFileSet() *fileSet {
	r := fileSet{
		files:		make([]*file, 0),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// FileSetIfc
// -------------------------------------------------------------------------------------------------

func (r *fileSet) AddFile(path string) {
	r.files = append(r.files, NewFile(path))
}

func (r *fileSet) Len() int {
	return len(r.files)
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

func (r *fileSet) GetIterator() func () interface{} {
        idx := 0
        var data_len = r.Len()
        return func () interface{} {
                // If we're done iterating, return nothing
                if idx >= data_len { return nil }
                prev_idx := idx
                idx++
                return r.files[prev_idx]
        }
}

