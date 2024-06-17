package uid

import "strconv"

// Nil nil
const Nil = Id(0)

// Id represents a global unique id.
type Id int64

// IsNil checks if an Id is nil.
func (id Id) IsNil() bool {
	return id == Nil
}

// String implements fmt.Stringer
func (id Id) String() string {
	return strconv.Itoa(int(id))
}
