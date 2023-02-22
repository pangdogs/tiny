package container

// NewAllocator 创建链表内存分配器
func NewAllocator[T any](size int) Allocator[T] {
	if size <= 0 {
		panic("size less equal 0 is invalid")
	}
	return &_Allocator[T]{
		size: size,
	}
}

// Allocator 链表内存分配器
type Allocator[T any] interface {
	// Alloc 分配链表元素
	Alloc() *Element[T]
}

type _Allocator[T any] struct {
	heap  []Element[T]
	index int
	size  int
}

// Alloc 分配链表元素
func (cache *_Allocator[T]) Alloc() *Element[T] {
	if cache.index >= len(cache.heap) {
		cache.index = 0
		cache.heap = make([]Element[T], cache.size)
	}

	e := &cache.heap[cache.index]
	cache.index++

	return e
}
