package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type Key interface {
	Less(than Key) bool
}

type keys []Key
type values []interface{}
type children []*node

type node struct {
	keys     keys
	children children
	values   values
}

type BplusTree struct {
	deep   int
	length int
	root   *node
}

func NewBplusTree(dp int) *BplusTree {
	return &BplusTree{
		deep: dp,
	}
}

func (b *BplusTree) Set(key Key, val interface{}) interface{} {
	if key == nil {
		panic(fmt.Errorf("nil key"))
	}
	if b.root == nil {
		b.root = &node{}
	}

	b.root.set(key, val, b.deep)

	if len(b.root.keys) >= b.deep {
		b.root = &node{
			children: []*node{b.root},
		}
		b.root.splitAndGrow(0, b.deep)
	}
	return nil
}

func (b *BplusTree) Get(key Key) interface{} {

	ptr, idx, found := b.root.get(key)

	if found {
		return ptr.values[idx]
	}
	return nil
}

func (b *BplusTree) Del(key Key) {

	b.root.del(key, b.deep)
	if len(b.root.keys) == 0 && len(b.root.children) == 1 {
		b.root = b.root.children[0]
	}
}

func (b *BplusTree) Iter() {

}

//------------------------------------------------------------

func (n *node) index(key Key) (int, bool) {
	if len(n.keys) == 0 {
		return 0, false
	}

	var idx int
	for idx = 0; idx < len(n.keys); idx++ {
		if key.Less(n.keys[idx]) {
			break
		}
	}

	if idx > 0 && !n.keys[idx-1].Less(key) {
		// a.Less(b) && b.Less(a) equal to  a == b
		return idx - 1, true
	}

	return idx, false
}

func (n *node) insertK(idx int, key Key) {

	n.keys = append(n.keys, key)
	if idx < len(n.keys) {
		copy(n.keys[idx+1:], n.keys[idx:])
	}
	n.keys[idx] = key
}

func (n *node) insertV(idx int, val interface{}) {

	n.values = append(n.values, val)
	if idx < len(n.values) {
		copy(n.values[idx+1:], n.values[idx:])
	}
	n.values[idx] = val
}

func (n *node) insertN(idx int, nd *node) {

	n.children = append(n.children, nil)
	if idx < len(n.children) {
		copy(n.children[idx+1:], n.children[idx:])
	}
	n.children[idx] = nd
}

func (n *node) deleteK(idx int) {

	if idx < len(n.keys)-1 {
		copy(n.keys[idx:], n.keys[idx+1:])
	}
	n.keys = n.keys[:len(n.keys)-1]
}

func (n *node) deleteV(idx int) {
	if idx < len(n.values)-1 {
		copy(n.values[idx:], n.values[idx+1:])
	}
	n.values = n.values[:len(n.values)-1]
}

func (n *node) deleteN(idx int) {
	if idx < len(n.children)-1 {
		copy(n.children[idx:], n.children[idx+1:])
	}
	n.children = n.children[:len(n.children)-1]
}

func (n *node) split(idx int) (right *node, k Key) {

	right = &node{
		keys: make(keys, len(n.keys)-idx),
	}
	copy(right.keys, n.keys[idx:])
	n.keys = n.keys[0:idx]

	if len(n.children) == 0 {
		// leaf
		right.values = make(values, len(n.values)-idx)
		copy(right.values, n.values[idx:])
		n.values = n.values[0:idx]

		return right, n.keys[idx-1]
	}
	// node
	right.children = make(children, len(n.children)-idx)
	copy(right.children, n.children[idx:])
	n.children = n.children[0:idx]

	k = n.keys[idx-1]
	n.keys = n.keys[0 : idx-1]

	return right, k
}

func ceil(dp int) int {
	return int(math.Ceil(float64(dp) / 2))
}

func floor(dp int) int {
	return int(math.Floor(float64(dp) / 2))
}

func (n *node) splitAndGrow(idx, dp int) (split bool) {

	nc := n.children[idx]
	if len(nc.keys) >= dp {

		right, k := nc.split(ceil(dp))
		n.insertK(idx, k)
		n.insertN(idx+1, right)
		return true
	}
	return false
}

func (n *node) set(key Key, val interface{}, dp int) {

	idx, found := n.index(key)
	if len(n.children) == 0 {
		// leaft
		if found {
			// replace
			n.values[idx] = val
		} else {
			// insert here
			n.insertK(idx, key)
			n.insertV(idx, val)
		}
		return
	}

	n.children[idx].set(key, val, dp)
	n.splitAndGrow(idx, dp)
	return
}

func (n *node) get(key Key) (ptr *node, idx int, found bool) {
	if n == nil {
		return nil, 0, false
	}

	idx, found = n.index(key)

	if found && len(n.children) == 0 {
		return n, idx, true
	}

	if len(n.children) == 0 {
		return n, idx, false
	}

	return n.children[idx].get(key)
}

func (n *node) del(key Key, dp int) {

	idx, found := n.index(key)
	if len(n.children) == 0 {
		if found {
			// delete
			n.deleteK(idx)
			n.deleteV(idx)
		}
		return
	}
	n.children[idx].del(key, dp)
	n.underFlow(idx, floor(dp))
}

func (n *node) underFlow(idx, minItem int) {
	nc := n.children[idx]

	if len(nc.children) == 0 && len(nc.keys) < minItem {
		// rotate or merge
		if idx > 0 && len(n.children[idx-1].keys) > minItem {
			// steal from left
			left := n.children[idx-1]
			key := left.keys[len(left.keys)-1]
			val := left.values[len(left.values)-1]
			sep := left.keys[len(left.keys)-2]
			left.deleteK(len(left.keys) - 1)
			left.deleteV(len(left.values) - 1)
			nc.insertK(0, key)
			nc.insertV(0, val)
			n.keys[idx-1] = sep
			return
		} else if idx < len(n.children)-1 && len(n.children[idx+1].keys) > minItem {
			// steal from rigt
			right := n.children[idx+1]
			key := right.keys[0]
			val := right.values[0]
			sep := key
			right.deleteK(0)
			right.deleteV(0)
			nc.insertK(len(nc.keys), key)
			nc.insertV(len(nc.values), val)
			n.keys[idx] = sep
			return
		} else {
			// merge adn remove sep
			if idx > 0 {
				left := n.children[idx-1]
				left.keys = append(left.keys, nc.keys...)
				left.values = append(left.values, nc.values)
				n.deleteN(idx)
				n.deleteK(idx - 1)
				return
			} else if idx < len(n.children)-1 {
				right := n.children[idx+1]
				nc.keys = append(nc.keys, right.keys...)
				nc.values = append(nc.values, right.values...)
				n.deleteK(idx)
				n.deleteN(idx + 1)
				return
			} else {
				fmt.Println("nil")
			}
		}
	} else {
		// rotate or merge
		if len(nc.children) != 0 && len(nc.children) < minItem {

			if idx > 0 && len(n.children[idx-1].keys) > minItem {
				// steal from left
				left := n.children[idx-1]
				key := left.keys[len(left.keys)-1]
				child := left.children[len(left.children)-1]
				sep := left.keys[len(left.keys)-2]
				left.deleteK(len(left.keys) - 1)
				left.deleteN(len(left.values) - 1)
				nc.insertK(0, key)
				nc.insertN(0, child)
				n.keys[idx] = sep
				return
			} else if idx < len(n.children)-1 && len(n.children[idx+1].keys) > minItem {
				// steal from rigt
				right := n.children[idx+1]
				key := right.keys[0]
				child := right.children[0]
				sep := key
				right.deleteK(0)
				right.deleteN(0)
				nc.insertK(len(nc.keys), key)
				nc.insertN(len(nc.children), child)
				n.keys[idx] = sep
				return
			} else {
				// merge adn remove sep
				if idx > 0 {
					// left
					left := n.children[idx-1]
					k := n.keys[idx-1]
					left.keys = append(left.keys, nc.keys...)
					left.keys = append(left.keys, k)
					left.children = append(left.children, nc.children...)
					n.deleteN(idx)
					n.deleteK(idx - 1)
					return
				} else if idx < len(n.children)-1 {
					// right
					right := n.children[idx+1]
					k := n.keys[idx]
					nc.keys = append(nc.keys, k)
					nc.keys = append(nc.keys, right.keys...)
					nc.children = append(nc.children, right.children...)
					n.deleteK(idx)
					n.deleteN(idx + 1)
					return
				} else {
					fmt.Println("nil")
				}
			}
		}
	}
}

// -----------------------------------------------
func (b *BplusTree) Print() {

	fmt.Println("Node --")
	defer fmt.Println("Node ++")
	lrtr := &node{}
	sep := &node{}

	queue := make([]*node, 0)
	queue = append(queue, b.root)
	queue = append(queue, lrtr)

	for len(queue) != 0 {

		for i := 0; i < len(queue); i++ {
			if queue[i] == nil {

			} else if queue[i] == lrtr {
				queue = queue[i+1 : len(queue)]
				queue = append(queue, lrtr)
				if len(queue) == 1 {
					fmt.Println(" ")
					return
				}
				fmt.Println(" ")
				break
			} else if queue[i] == sep {
				fmt.Print("|")
			} else {
				fmt.Print(queue[i].keys, " ")
				queue = append(queue, queue[i].children...)
				queue = append(queue, sep)
			}
		}
	}
}
func (n *node) dot() string {
	tmpl := ""
	var s []string
	for _, v := range n.keys {
		s = append(s, fmt.Sprintf("%v", v))
	}
	tmpl += fmt.Sprintf("\t\"%p\" [label=\"%v\"];\n", n, strings.Join(s, "-"))

	if len(n.children) == 0 {

	} else {
		for _, v := range n.children {
			tmpl += fmt.Sprintf("\t\"%p\" ->  \"%p\" ;\n", n, v)
			tmpl += v.dot()
		}
	}

	return tmpl
}

func (b *BplusTree) Dot(name string) {

	tpl := " digraph graphname {\n"
	tpl += b.root.dot()
	tpl += "}"

	ioutil.WriteFile(name, []byte(tpl), os.FileMode(0666))
}
