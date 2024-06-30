package generic

// NewList 创建链表
func NewList[T any]() *List[T] {
	return &List[T]{}
}

// List 链表，可以在遍历时在任意位置添加或删除元素，递归添加或删除元素时仍然能正常工作，非线程安全。
type List[T any] struct {
	New  Func1[T, *Node[T]]
	root Node[T]
	len  int
	ver  int64
}

// Len 链表长度
func (l *List[T]) Len() int {
	return l.len
}

// Version 数据变化版本
func (l *List[T]) Version() int64 {
	return l.ver
}

// Front 链表头部
func (l *List[T]) Front() *Node[T] {
	if l.len <= 0 {
		return nil
	}
	return l.root._next
}

// Back 链表尾部
func (l *List[T]) Back() *Node[T] {
	if l.len <= 0 {
		return nil
	}
	return l.root._prev
}

// Remove 删除元素
func (l *List[T]) Remove(n *Node[T]) T {
	if n.list == l {
		n.Escape()
	}
	return n.V
}

// PushFront 在链表头部插入数据
func (l *List[T]) PushFront(value T) *Node[T] {
	l.lazyInit()
	return l.insertValue(value, &l.root)
}

// PushBack 在链表尾部插入数据
func (l *List[T]) PushBack(value T) *Node[T] {
	l.lazyInit()
	return l.insertValue(value, l.root._prev)
}

// InsertBefore 在链表指定位置前插入数据
func (l *List[T]) InsertBefore(value T, at *Node[T]) *Node[T] {
	if !l.check(at) {
		return nil
	}
	return l.insertValue(value, at._prev)
}

// InsertAfter 在链表指定位置后插入数据
func (l *List[T]) InsertAfter(value T, at *Node[T]) *Node[T] {
	if !l.check(at) {
		return nil
	}
	return l.insertValue(value, at)
}

// MoveToFront 移动元素至链表头部
func (l *List[T]) MoveToFront(n *Node[T]) {
	if !l.check(n) || l.root._next == n {
		return
	}
	l.move(n, &l.root)
}

// MoveToBack 移动元素至链表尾部
func (l *List[T]) MoveToBack(n *Node[T]) {
	if !l.check(n) || l.root._prev == n {
		return
	}
	l.move(n, l.root._prev)
}

// MoveBefore 移动元素至链表指定位置前
func (l *List[T]) MoveBefore(n, at *Node[T]) {
	if !l.check(n) || !l.check(at) || n == at {
		return
	}
	l.move(n, at._prev)
}

// MoveAfter 移动元素至链表指定位置后
func (l *List[T]) MoveAfter(n, at *Node[T]) {
	if !l.check(n) || !l.check(at) || n == at {
		return
	}
	l.move(n, at)
}

// PushFrontList 在链表头部插入其他链表，可以传入自身
func (l *List[T]) PushFrontList(other *List[T]) {
	if other == nil {
		return
	}
	l.lazyInit()
	for i, n := other.Len(), other.Back(); i > 0; i, n = i-1, n.Prev() {
		l.insertValue(n.V, &l.root)
	}
}

// PushBackList 在链表尾部插入其他链表，可以传入自身
func (l *List[T]) PushBackList(other *List[T]) {
	if other == nil {
		return
	}
	l.lazyInit()
	for i, n := other.Len(), other.Front(); i > 0; i, n = i-1, n.Next() {
		l.insertValue(n.V, l.root._prev)
	}
}

// Traversal 遍历元素
func (l *List[T]) Traversal(visitor func(n *Node[T]) bool) {
	if visitor == nil {
		return
	}

	for n := l.Front(); n != nil; n = n.Next() {
		if !visitor(n) {
			break
		}
	}
}

// TraversalAt 从指定位置开始遍历元素
func (l *List[T]) TraversalAt(visitor func(n *Node[T]) bool, at *Node[T]) {
	if visitor == nil || !l.check(at) {
		return
	}

	for n := at; n != nil; n = n.Next() {
		if !visitor(n) {
			break
		}
	}
}

// ReversedTraversal 反向遍历元素
func (l *List[T]) ReversedTraversal(visitor func(n *Node[T]) bool) {
	if visitor == nil {
		return
	}

	for n := l.Back(); n != nil; n = n.Prev() {
		if !visitor(n) {
			break
		}
	}
}

// ReversedTraversalAt 从指定位置开始反向遍历元素
func (l *List[T]) ReversedTraversalAt(visitor func(n *Node[T]) bool, at *Node[T]) {
	if visitor == nil || !l.check(at) {
		return
	}

	for n := at; n != nil; n = n.Prev() {
		if !visitor(n) {
			break
		}
	}
}

// lazyInit 延迟初始化
func (l *List[T]) lazyInit() {
	if l.root._next != nil {
		return
	}
	l.root._next = &l.root
	l.root._prev = &l.root
	if l.New == nil {
		l.New = newNode[T]
	}
}

// insertValue 插入数据
func (l *List[T]) insertValue(value T, at *Node[T]) *Node[T] {
	l.lazyInit()
	return l.insert(l.New(value), at)
}

// insert 插入元素
func (l *List[T]) insert(n, at *Node[T]) *Node[T] {
	n._prev = at
	n._next = at._next
	n._prev._next = n
	n._next._prev = n
	n.list = l
	l.len++
	l.ver++
	n.ver = l.ver
	return n
}

// remove 删除元素
func (l *List[T]) remove(n *Node[T]) {
	l.lazyInit()
	n._prev._next = n._next
	n._next._prev = n._prev
	l.len--
	l.ver++
}

// move 移动元素
func (l *List[T]) move(n, at *Node[T]) *Node[T] {
	l.lazyInit()

	if n == at {
		return n
	}
	n._prev._next = n._next
	n._next._prev = n._prev

	n._prev = at
	n._next = at._next
	n._prev._next = n
	n._next._prev = n

	l.ver++

	return n
}

func (l *List[T]) check(n *Node[T]) bool {
	return n != nil && n.list == l
}
