package generic

func newNode[T any](value T) *Node[T] {
	return &Node[T]{
		V: value,
	}
}

// Node 链表节点
type Node[T any] struct {
	V            T
	_next, _prev *Node[T]
	list         *List[T]
	ver          int64
}

// Version 创建时的数据变化版本
func (n *Node[T]) Version() int64 {
	return n.ver
}

// Next 下一个元素
func (n *Node[T]) Next() *Node[T] {
	for n := n.next(); n != nil; n = n.next() {
		if !n.Escaped() {
			return n
		}
	}
	return nil
}

// Prev 前一个元素
func (n *Node[T]) Prev() *Node[T] {
	for p := n.prev(); p != nil; p = p.prev() {
		if !p.Escaped() {
			return p
		}
	}
	return nil
}

// Escape 从链表中删除
func (n *Node[T]) Escape() {
	if n.list != nil {
		n.list.remove(n)
		n.list = nil
	}
}

// Escaped 是否已从链表中删除
func (n *Node[T]) Escaped() bool {
	return n.list == nil
}

// next 下一个元素，包含正在删除的元素
func (n *Node[T]) next() *Node[T] {
	if n := n._next; n.list != nil && n != &n.list.root {
		return n
	}
	return nil
}

// prev 前一个元素，包含正在删除的元素
func (n *Node[T]) prev() *Node[T] {
	if p := n._prev; n.list != nil && p != &n.list.root {
		return p
	}
	return nil
}
