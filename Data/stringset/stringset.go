package stringset

/*
This utility library helps to deal with common scenarios involving sets of unique strings!
*/

type StringSetIfc interface {
	Set(s string)
	SetAll(ss *[]string)
	Merge(stringSet StringSetIfc)

	Drop(s string)
	DropAll(ss *[]string)

	IsEmpty() bool
	Size() int
	Has(s string) bool
	HasAll(ss *[]string) bool
	HasAny(ss *[]string) bool

	ToArray() *[]string

	GetIterator() func () interface{}
}

type StringSet struct {
	stringset		map[string]bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewStringSet() *StringSet {
	return &StringSet{
		stringset:		make(map[string]bool),
	}
}

// -------------------------------------------------------------------------------------------------
// StringSetIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *StringSet) Set(s string) {
	r.stringset[s] = true
}

func (r *StringSet) SetAll(ss *[]string) {
	if nil == ss { return }
	for _, s := range *ss { r.Set(s) }
}

func (r *StringSet) Merge(stringSet StringSetIfc) {
	if nil == stringSet { return }
	r.SetAll(stringSet.ToArray())
}

func (r *StringSet) Drop(s string) {
	delete(r.stringset, s)
}

func (r *StringSet) DropAll(ss *[]string) {
	if nil == ss { return }
	for _, s := range *ss { r.Drop(s) }
}

func (r *StringSet) IsEmpty() bool {
	return r.Size() == 0
}

func (r *StringSet) Size() int {
	return len(r.stringset)
}

func (r *StringSet) Has(s string) bool {
	_, ok := r.stringset[s]
	return ok
}

func (r *StringSet) HasAll(ss *[]string) bool {
	for _, s := range *ss {
		if ! r.Has(s) { return false }
	}
	return true
}

func (r *StringSet) HasAny(ss *[]string) bool {
	for _, s := range *ss {
		if r.Has(s) { return true }
	}
	return false
}

func (r *StringSet) ToArray() *[]string {
	ss := make([]string, r.Size())
	i := 0
	for s, _ := range r.stringset {
		ss[i] = s
		i++
	}
	return &ss
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *StringSet) GetIterator() func () interface{} {
	idx := 0
	var data_len = r.Size()
	ss := r.ToArray()
	return func () interface{} {
		// If we're done iterating, return do nothing
		if idx >= data_len { return nil }
		prev_idx := idx
		idx++
		return &((*ss)[prev_idx])
	}

	return nil
}
