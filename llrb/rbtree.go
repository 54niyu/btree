package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Key interface {
	Less(than Key) bool
}

type rbNode struct {
	left, right *rbNode
	key         Key
	val         interface{}
	red         bool
}

type RBTree struct {
	root *rbNode
}

func NewRBTree() *RBTree {
	return new(RBTree)
}

func (t *RBTree) Set(key Key, val interface{}) {
	if key == nil {
		panic(fmt.Errorf("key is nil"))
	}
	t.root = t.root.set(key, val)
	if t.root.red {
		t.root.red = false
	}
	return
}

func (t *RBTree) Get(key Key) (val interface{}, found bool) {
	if key == nil {
		panic(fmt.Errorf("key is nil"))
	}
	if t.root == nil {
		return nil, false
	}
	return t.root.get(key)
}

func (t *RBTree) Del(key Key) {
	if key == nil {
		panic(fmt.Errorf("key is nil"))
	}

	if t.root == nil {
		return
	}

	t.root = t.root.del(key)
	if t.root != nil {
		t.root.red = false
	}
	return
}

type Iterator func(k Key, val interface{}) error

func (t *RBTree) Iter(f Iterator) error {
	return t.root.iterator(f)
}

func (t *RBTree) Hight() {
	hight(t.root, 0)
}

func hight(n *rbNode, h int) {
	if n == nil {
		fmt.Print(h, " ")
		return
	}
	if !n.red {
		h += 1
	}
	hight(n.left, h)
	hight(n.right, h)
}

//-----------------------------------------------------------------------------

func (n *rbNode) iterator(f Iterator) error {
	if n == nil {
		return nil
	}
	if err := n.left.iterator(f); err != nil {
		return err
	}
	if err := f(n.key, n.val); err != nil {
		return err
	}
	if err := n.right.iterator(f); err != nil {
		return err
	}
	return nil
}

func (n *rbNode) compare(key Key) int {
	if key.Less(n.key) {
		return -1
	}
	if !n.key.Less(key) {
		return 0
	} else {
		return 1
	}
}

func (n *rbNode) isRed() bool {
	if n == nil {
		return false
	}
	return n.red
}

func (n *rbNode) rotateRight() *rbNode {

	if !n.left.red {
		panic("rotateRight left is black")
	}

	x := n.left
	n.left = x.right
	x.right = n
	x.red = x.right.red
	x.right.red = true

	return x
}

func (n *rbNode) rotateLeft() *rbNode {

	if !n.right.red {
		panic("rotateLeft right is black")
	}

	x := n.right
	n.right = x.left
	x.left = n
	x.red = x.left.red
	x.left.red = true

	return x
}

func (n *rbNode) colorFlip() {
	n.red = !n.red
	n.left.red = !n.left.red
	n.right.red = !n.right.red
}

func (n *rbNode) get(key Key) (val interface{}, found bool) {
	if n == nil {
		return nil, false
	}
	if key.Less(n.key) {
		return n.left.get(key)
	} else if n.key.Less(key) {
		return n.right.get(key)
	} else {
		return n.val, true
	}
}

func (n *rbNode) set(key Key, val interface{}) *rbNode {
	if n == nil {
		return &rbNode{key: key, val: val, red: true}
	}

	/*
		2-3-4 tree
		if n.left.isRed() && n.right.isRed() {
			n.colorFlip()
		}
	*/

	c := n.compare(key)
	switch c {
	case 0:
		n.val = val
	case -1:
		n.left = n.left.set(key, val)
	case 1:
		n.right = n.right.set(key, val)
	}

	if n.right.isRed() {
		n = n.rotateLeft()
	}

	if n.left.isRed() && n.left.left.isRed() {
		n = n.rotateRight()
	}
	/*
		2-3 tree
	*/
	if n.left.isRed() && n.right.isRed() {
		n.colorFlip()
	}

	return n
}

func (n *rbNode) del(key Key) *rbNode {
	if n == nil {
		return nil
	}

	if key.Less(n.key) {
		if n.left != nil && !n.left.isRed() && !n.left.left.isRed() {
			n = n.moveRedLeft()
		}
		n.left = n.left.del(key)
	} else {
		if n.left.isRed() {
			n = n.rotateRight()
		}
		if !n.key.Less(key) && n.right == nil {
			return nil
		}
		if n.right != nil && !n.right.isRed() && !n.right.left.isRed() {
			n = n.moveRedRight()
		}
		if !n.key.Less(key) {
			v := n.right.minInRight()
			n.key = v.key
			n.val = v.val
			n.right = n.right.delMin()
		} else {
			n.right = n.right.del(key)
		}
	}

	n = n.fixUp()
	return n
}

func (n *rbNode) minInRight() *rbNode {
	ptr := n
	for {
		if ptr.left != nil {
			ptr = ptr.left
		} else {
			break
		}
	}
	return ptr
}

func (n *rbNode) maxInLeft() *rbNode {
	ptr := n
	for {
		if ptr.right != nil {
			ptr = ptr.right
		} else {
			break
		}
	}
	return ptr
}

func (n *rbNode) leanRight() {}

func (n *rbNode) delMin() *rbNode {
	if n == nil {
		return nil
	}

	if n.left == nil {
		return nil
	}

	if n.left != nil && !n.left.isRed() && !n.left.left.isRed() {
		n = n.moveRedLeft()
	}

	n.left = n.left.delMin()
	n = n.fixUp()
	return n
}

func (n *rbNode) delMax() *rbNode {
	if n == nil {
		return nil
	}

	if n.left.isRed() {
		n = n.rotateRight()
	}

	if n.right == nil {
		return nil
	}

	if n.right != nil && !n.right.isRed() && !n.right.left.isRed() {
		n = n.moveRedRight()
	}

	n.right = n.right.delMax()
	n = n.fixUp()
	return n
}

func (n *rbNode) moveRedLeft() *rbNode {
	n.colorFlip()
	if n.right != nil && n.right.left.isRed() {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.colorFlip()
	}
	return n
}

func (n *rbNode) moveRedRight() *rbNode {
	n.colorFlip()
	if n.left != nil && n.left.left.isRed() {
		n = n.rotateRight()
		n.colorFlip()
	}
	return n
}

func (n *rbNode) fixUp() *rbNode {
	if n.right.isRed() {
		n = n.rotateLeft()
	}
	if n.left.isRed() && n.left.left.isRed() {
		n = n.rotateRight()
	}
	if n.left.isRed() && n.right.isRed() {
		n.colorFlip()
	}
	return n
}

func (t *RBTree) I() {
	t.Iter(Iterator(func(key Key, val interface{}) error {
		fmt.Print(key, " ")
		return nil
	}))
}

func (n *rbNode) dot() string {
	if n == nil {
		return ""
	}

	tmpl := ""
	tmpl += fmt.Sprintf("\t\"%p\" [label=\"<f0> | <f1> %v | <f2>\" color=%v shape=record];\n", n, n.key, func() string {
		if n.red {
			return "red"
		} else {
			return "black"
		}
	}())

	if n.left != nil {
		tmpl += fmt.Sprintf("\t\"%p\":f0 ->  \"%p\":f1 [color=blue style=dotted] ;\n", n, n.left)
		tmpl += n.left.dot()
	}
	if n.right != nil {
		tmpl += fmt.Sprintf("\t\"%p\":f2 ->  \"%p\":f1 ;\n", n, n.right)
		tmpl += n.right.dot()
	}

	return tmpl
}

func (b *RBTree) Dot(name string) {

	tpl := " digraph graphname {\n"
	tpl += b.root.dot()
	tpl += "}"

	ioutil.WriteFile(name, []byte(tpl), os.FileMode(0666))
}
