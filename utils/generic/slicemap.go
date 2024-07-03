package generic

import (
	"git.golaxy.org/tiny/utils/types"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func MakeSliceMap[K constraints.Ordered, V any](kvs ...KV[K, V]) SliceMap[K, V] {
	m := make(SliceMap[K, V], 0, len(kvs))
	for i := range kvs {
		kv := kvs[i]
		m.Add(kv.K, kv.V)
	}
	return m
}

func NewSliceMap[K constraints.Ordered, V any](kvs ...KV[K, V]) *SliceMap[K, V] {
	m := make(SliceMap[K, V], 0, len(kvs))
	for i := range kvs {
		kv := kvs[i]
		m.Add(kv.K, kv.V)
	}
	return &m
}

func MakeSliceMapFromGoMap[K constraints.Ordered, V any](m map[K]V) SliceMap[K, V] {
	sm := make(SliceMap[K, V], 0, len(m))
	for k, v := range m {
		sm.Add(k, v)
	}
	return sm
}

func NewSliceMapFromGoMap[K constraints.Ordered, V any](m map[K]V) *SliceMap[K, V] {
	sm := make(SliceMap[K, V], 0, len(m))
	for k, v := range m {
		sm.Add(k, v)
	}
	return &sm
}

type KV[K constraints.Ordered, V any] struct {
	K K
	V V
}

type SliceMap[K constraints.Ordered, V any] []KV[K, V]

func (m *SliceMap[K, V]) Add(k K, v V) {
	idx, ok := slices.BinarySearchFunc(*m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	if ok {
		(*m)[idx] = KV[K, V]{K: k, V: v}
	} else {
		*m = slices.Insert(*m, idx, KV[K, V]{K: k, V: v})
	}
}

func (m *SliceMap[K, V]) TryAdd(k K, v V) bool {
	idx, ok := slices.BinarySearchFunc(*m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	if !ok {
		*m = slices.Insert(*m, idx, KV[K, V]{K: k, V: v})
	}
	return !ok
}

func (m *SliceMap[K, V]) Delete(k K) bool {
	idx, ok := slices.BinarySearchFunc(*m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	if ok {
		*m = slices.Delete(*m, idx, idx+1)
	}
	return ok
}

func (m SliceMap[K, V]) Get(k K) (V, bool) {
	idx, ok := slices.BinarySearchFunc(m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	if ok {
		return m[idx].V, true
	}
	return types.ZeroT[V](), false
}

func (m SliceMap[K, V]) Value(k K) V {
	idx, ok := slices.BinarySearchFunc(m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	if ok {
		return m[idx].V
	}
	return types.ZeroT[V]()
}

func (m SliceMap[K, V]) Exist(k K) bool {
	_, ok := slices.BinarySearchFunc(m, KV[K, V]{K: k}, func(a, b KV[K, V]) int {
		return compare(a.K, b.K)
	})
	return ok
}

func (m SliceMap[K, V]) Len() int {
	return len(m)
}

func (m SliceMap[K, V]) Range(fun Func2[K, V, bool]) {
	for _, kv := range m {
		if !fun.Exec(kv.K, kv.V) {
			return
		}
	}
}

func (m SliceMap[K, V]) Each(fun Action2[K, V]) {
	for _, kv := range m {
		fun.Exec(kv.K, kv.V)
	}
}

func (m SliceMap[K, V]) ReversedRange(fun Func2[K, V, bool]) {
	for i := len(m) - 1; i >= 0; i-- {
		kv := m[i]
		if !fun.Exec(kv.K, kv.V) {
			return
		}
	}
}

func (m SliceMap[K, V]) ReversedEach(fun Action2[K, V]) {
	for i := len(m) - 1; i >= 0; i-- {
		kv := m[i]
		fun.Exec(kv.K, kv.V)
	}
}

func (m SliceMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.Len())
	for _, kv := range m {
		keys = append(keys, kv.K)
	}
	return keys
}

func (m SliceMap[K, V]) Values() []V {
	values := make([]V, 0, m.Len())
	for _, kv := range m {
		values = append(values, kv.V)
	}
	return values
}

func (m SliceMap[K, V]) Clone() SliceMap[K, V] {
	return slices.Clone(m)
}

func (m SliceMap[K, V]) ToUnorderedSliceMap() UnorderedSliceMap[K, V] {
	rv := make(UnorderedSliceMap[K, V], 0, len(m))
	for _, kv := range m {
		rv.Add(kv.K, kv.V)
	}
	return rv
}

func (m SliceMap[K, V]) ToGoMap() map[K]V {
	rv := make(map[K]V, len(m))
	for _, kv := range m {
		rv[kv.K] = kv.V
	}
	return rv
}

func compare[T constraints.Ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}
