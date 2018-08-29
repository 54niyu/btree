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
	fmt.Println("Del ", key)

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

func (n *rbNode) rotateRight() {

	l := n.left.left
	m := n.left.right
	r := n.right

	n.left.right = r
	n.left.left = m

	n.right = n.left
	n.left = l

	t, t2 := n.key, n.val
	n.key, n.val = n.right.key, n.right.val
	n.right.key, n.right.val = t, t2
}

func (n *rbNode) rotateLeft() {

	l := n.left
	m := n.right.left
	r := n.right.right

	n.right.left = l
	n.right.right = m

	n.left = n.right
	n.right = r

	t, t2 := n.key, n.val
	n.key, n.val = n.left.key, n.left.val
	n.left.key, n.left.val = t, t2
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

	if n.left.isRed() && n.right.isRed() {
		n.colorFlip()
	}

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
		n.rotateLeft()
	}

	if n.left.isRed() && n.left.left.isRed() {
		n.rotateRight()
	}

	return n
}

func (n *rbNode) del(key Key) *rbNode {
	if n == nil {
		return nil
	}

	c := n.compare(key)
	if c == -1 {
		if n.left != nil && !n.left.isRed() && !n.left.left.isRed() {
			n.moveRedLeft()
		}
		n.left = n.left.del(key)
	} else {
		if n.left.isRed() {
			n.rotateRight()
			c = n.compare(key)
		}
		if c == 0 && n.right == nil {
			return nil
		}
		if n.right != nil && !n.right.isRed() && !n.right.left.isRed() {
			n.moveRedRight()
			c = n.compare(key)
		}
		if c == 0 {
			v := n.right.minInRight()
			n.key = v.key
			n.val = v.val
			n.right = n.right.delMin()
		} else {
			n.right = n.right.del(key)
		}
	}

	n.fixUp()
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
		n.moveRedLeft()
	}

	n.left = n.left.delMin()
	n.fixUp()
	return n
}

func (n *rbNode) delMax() {
	if n == nil {
		return
	}

	if n.left.isRed() {
		n.rotateRight()
	}

	if n.right == nil {
		return
	}

	if n.right != nil && !n.right.isRed() && !n.right.left.isRed() {
		n.moveRedRight()
	}

	n.right.delMax()
	n.fixUp()
	return
}

func (n *rbNode) moveRedLeft() {
	n.colorFlip()
	if n.right != nil && n.right.left.isRed() {
		n.right.rotateRight()
		n.rotateLeft()
		n.colorFlip()
	}
}

func (n *rbNode) moveRedRight() {
	n.colorFlip()
	if n.left != nil && n.left.left.isRed() {
		n.rotateRight()
		n.colorFlip()
	}
}

func (n *rbNode) fixUp() {
	if n.right.isRed() {
		n.rotateLeft()
	}
	if n.left.isRed() && n.left.left.isRed() {
		n.rotateRight()
	}
	if n.left.isRed() && n.right.isRed() {
		n.colorFlip()
	}
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
