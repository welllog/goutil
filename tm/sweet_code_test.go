package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

type Person1 struct {
	a bool
	b int64
	c int8
	d string
}

func TestUnsafe(t *testing.T) {
	a := Person1{}
	fmt.Println(unsafe.Offsetof(a.a))
	fmt.Println(unsafe.Offsetof(a.b))
	fmt.Println(unsafe.Offsetof(a.c))
	fmt.Println(unsafe.Offsetof(a.d))
	fmt.Println("--------")
	fmt.Println(unsafe.Alignof(a.a))
	fmt.Println(unsafe.Alignof(a.b))
	fmt.Println(unsafe.Alignof(a.c))
	fmt.Println(unsafe.Alignof(a.d))
	fmt.Println("--------")
	fmt.Println(unsafe.Sizeof(a.a))
	fmt.Println(unsafe.Sizeof(a.b))
	fmt.Println(unsafe.Sizeof(a.c))
	fmt.Println(unsafe.Sizeof(a.d))

	p := unsafe.Pointer(&a)
	up0 := uintptr(p)
	up := up0 + unsafe.Offsetof(a.b)
	p = unsafe.Pointer(up)
	pb := (*int64)(p)
	*pb = 2
	fmt.Println(a)

	a.b = 8
	abp := unsafe.Pointer(&a.b)
	fmt.Println(*(*int64)(atomic.LoadPointer(&abp)))
}

type Node struct {
	Data  interface{}
	Left  *Node
	Right *Node
}

type NodeStack struct {
	stack []*Node
}

func NewNodeStack(cap int) *NodeStack {
	return &NodeStack{stack: make([]*Node, 0, cap)}
}

func (stack *NodeStack) Push(node *Node) {
	stack.stack = append(stack.stack, node)
}

func (stack *NodeStack) Pop() *Node {
	l := len(stack.stack)
	if l == 0 {
		return nil
	}
	node := stack.stack[l-1]
	stack.stack = stack.stack[:l-1]
	return node
}

func (stack *NodeStack) Top() *Node {
	l := len(stack.stack)
	if l == 0 {
		return nil
	}
	return stack.stack[l-1]
}

func (stack *NodeStack) IsEmpty() bool {
	return len(stack.stack) == 0
}

// 广度优先遍历
func (n *Node) BreadthFirst() {
	queue := make([]*Node, 0, 40)
	queue = append(queue, n)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Print(node.Data, " ")
		if node.Left != nil {
			queue = append(queue, node.Left)
		}

		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}

// 前序遍历
func (n *Node) PreOrder() {
	stack := NewNodeStack(40)
	stack.Push(n)
	for !stack.IsEmpty() {
		node := stack.Pop()
		fmt.Print(node.Data, " ")

		if node.Right != nil {
			stack.Push(node.Right)
		}

		if node.Left != nil {
			stack.Push(node.Left)
		}

	}
}

func (n *Node) PreOrder1() {
	stack := NewNodeStack(40)
	node := n
	for node != nil || !stack.IsEmpty() {
		for node != nil {
			fmt.Print(node.Data, " ")
			stack.Push(node)
			node = node.Left
		}
		// 回溯到父节点
		node = stack.Pop()
		node = node.Right
	}
}

// 中序遍历
func (n *Node) InOrder() {
	stack := NewNodeStack(40)
	node := n
	for node != nil || !stack.IsEmpty() {
		for node != nil {
			stack.Push(node)
			node = node.Left
		}
		// 回溯到父节点
		node = stack.Pop()
		fmt.Print(node.Data, " ")
		node = node.Right
	}
}

// 后序遍历
func (n *Node) TailOrder() {
	stack := NewNodeStack(40)
	node := n
	var cnode *Node
	for node != nil || !stack.IsEmpty() {
		for node != nil {
			stack.Push(node)
			node = node.Left
		}
		// 回溯到父节点
		node = stack.Top()
		if node.Right != nil && node.Right != cnode {
			node = node.Right
		} else {
			fmt.Print(node.Data, " ")
			stack.Pop()
			cnode = node
			node = nil
		}
	}
}

func TestTree(t *testing.T) {
	tree5 := &Node{
		Data:  5,
		Left:  &Node{Data: 7},
		Right: &Node{Data: 8},
	}
	tree2 := &Node{
		Data:  2,
		Left:  &Node{Data: 4},
		Right: tree5,
	}
	tree3 := &Node{
		Data:  3,
		Right: &Node{Data: 6},
	}
	tree1 := &Node{
		Data:  1,
		Left:  tree2,
		Right: tree3,
	}
	fmt.Print("广度: ")
	tree1.BreadthFirst()
	fmt.Println()
	fmt.Print("前序: ")
	tree1.PreOrder()
	fmt.Println()
	fmt.Print("前序1: ")
	tree1.PreOrder1()
	fmt.Println()
	fmt.Print("中序: ")
	tree1.InOrder()
	fmt.Println()
	fmt.Print("后序: ")
	tree1.TailOrder()
}
