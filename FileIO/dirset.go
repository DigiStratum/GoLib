package fileio

/*

A set of file system directories

Caller responsible for asserting that each entry is, could, or should be a directory.

*/

type DirSetIfc interface {
	AddDir(path string)
	Len() int
}

type dirSet struct {
	dirs []*dir
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewDirSet() *dirSet {
	r := dirSet{
		dirs: make([]*dir, 0),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// DirSetIfc
// -------------------------------------------------------------------------------------------------

func (r *dirSet) AddDir(path string) {
	r.dirs = append(r.dirs, Dir(path))
}

func (r *dirSet) Len() int {
	return len(r.dirs)
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

func (r *dirSet) GetIterator() func() interface{} {
	idx := 0
	var data_len = r.Len()
	return func() interface{} {
		// If we're done iterating, return nothing
		if idx >= data_len {
			return nil
		}
		prev_idx := idx
		idx++
		return r.dirs[prev_idx]
	}
}
