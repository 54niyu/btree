package main

import (
	"fmt"
	"math"
)

type Bnode struct {
	Val      []int
	Parent   *Bnode
	Children []*Bnode
	x, y     float64
}

type Btree struct {
	Deep int
	Root *Bnode
}

func NewBtree(dp int) *Btree {
	return &Btree{
		Deep: dp,
		Root: nil,
	}
}

func (b *Btree) Insert(v int) error {
	if b.Root == nil {
		b.Root = &Bnode{
			Parent:   nil,
			Val:      []int{v},
			Children: []*Bnode{nil, nil},
		}
		return nil
	}
	val, bubble, _ := b.Root.Insert(v, b.Deep)
	if bubble {
		b.Root = val
	}
	return nil
}

func (b *Bnode) InsertV(v, idx int) {
	b.Val = append(b.Val, 0)
	if idx < len(b.Val) {
		copy(b.Val[idx+1:], b.Val[idx:])
	}
	b.Val[idx] = v
}

func (b *Bnode) InsertN(n *Bnode, idx int) {
	b.Children = append(b.Children, nil)
	if idx < len(b.Children) {
		copy(b.Children[idx+1:], b.Children[idx:])
	}
	b.Children[idx] = n
}

func CopyNode(val []int, children []*Bnode) *Bnode {
	node := &Bnode{
		Val:      make([]int, len(val)),
		Children: make([]*Bnode, len(children)),
	}
	copy(node.Val[:], val[:])
	copy(node.Children[:], children[:])

	for _, v := range node.Children {
		if v != nil {
			v.Parent = node
		}
	}
	return node
}

func (b *Bnode) Insert(v, dp int) (*Bnode, bool, error) {
	if b == nil {
		return nil, false, nil
	}
	var ptr *Bnode
	var idx int
	for _, val := range b.Val {
		if v > val {
			idx++
		} else {
			break
		}
	}

	ptr = b.Children[idx]
	if ptr == nil {
		// 子节点 插入到当前位置
		b.InsertV(v, idx)
		b.InsertN(nil, idx)
	} else {
		// 向下深入
		val, bubble, _ := ptr.Insert(v, dp)
		if bubble {
			b.InsertV(val.Val[0], idx)
			b.Children[idx] = val.Children[0]
			b.InsertN(val.Children[1], idx+1)
			for _, v := range val.Children {
				v.Parent = b
			}
		}
	}

	if len(b.Children) > dp {
		// 分裂and上浮
		m := len(b.Val) / 2

		lf := b.Val[:m]
		md := b.Val[m]
		rt := b.Val[m+1:]
		clf := b.Children[:m+1]
		crt := b.Children[m+1:]

		blf := CopyNode(lf, clf)
		brt := CopyNode(rt, crt)

		bmd := &Bnode{
			Val:      []int{md},
			Children: []*Bnode{blf, brt},
		}
		brt.Parent = bmd
		blf.Parent = bmd

		return bmd, true, nil
	}

	return nil, false, nil
}

func (b *Btree) Get(v int) error {
	return b.Root.Get(v)
}

func (b *Bnode) Get(v int) error {
	ptr, idx, err := b.get(v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ptr != nil {
		fmt.Printf("Get %v %v \n", ptr.Val[idx], ptr.Val)
	}
	return nil
}

func (b *Bnode) get(v int) (*Bnode, int, error) {

	if b == nil {
		return nil, 0, nil
	}

	var (
		ptr   *Bnode
		idx   int
		found bool
	)

	for _, val := range b.Val {
		if v > val {
			idx++
		} else if val == v {
			idx++
			found = true
			break
		}
	}
	ptr = b.Children[idx]

	if found {
		return b, idx, nil
	}
	return ptr.get(v)
}

func (b *Btree) Delete(v int) error {
	return b.Root.Delete(v, b.Deep)
}

func (b *Bnode) Delete(v, deep int) error {
	fmt.Println("delete node ", v)
	defer fmt.Println("delete node over ")

	ptr, idx, _ := b.get(v)
	if ptr == nil {
		return nil
	}
	idx--

	var leaf bool
	for _, v := range ptr.Children {
		if v != nil {
			leaf = false
			break
		}
		leaf = true
	}

	if leaf {
		fmt.Println("leaf")

		ptr.Val = append(ptr.Val[:idx], ptr.Val[idx+1:]...)
		ptr.Children = append(ptr.Children[:idx], ptr.Children[idx+1:]...)

		min := int(math.Ceil(float64(deep) / 2.0))
		if len(ptr.Children) < min {
			ptr.Rebalance(min)
		}

	} else {
		fmt.Println("seperator")

		place, node := ptr.Children[idx].MaxInLeft()
		if node == nil {
			fmt.Println("nil")
		} else {
			// 替换并删除
			fmt.Printf("Replace %v for %v\n", place, node.Val)
			ptr.Val[idx] = node.Val[place]
			node.Val = node.Val[:len(node.Val)-1]
			node.Children = node.Children[:len(node.Children)-1]

			min := int(math.Ceil(float64(deep) / 2.0))
			if len(node.Children) < min {
				node.Rebalance(min)
			}
		}
	}

	return nil
}

func (b *Bnode) MaxInLeft() (idx int, node *Bnode) {

	if b == nil {
		return
	}
	rightChild := b.Children[len(b.Children)-1]

	if rightChild == nil {
		return len(b.Val) - 1, b
	}
	return rightChild.MaxInLeft()
}

func (b *Bnode) MaxInRight() (idx int, node *Bnode) {

	if b == nil {
		return
	}
	rightChild := b.Children[0]

	if rightChild == nil {
		return 0, b
	}
	return rightChild.MaxInRight()
}

func (b *Bnode) RotateLeft(idx int) {
	ridx := idx + 1
	l := b.Children[idx]
	r := b.Children[ridx]

	l.Val = append(l.Val, b.Val[idx])
	l.Children = append(l.Children, r.Children[0])
	b.Val[idx] = r.Val[0]
	r.Val = r.Val[1:]
	r.Children = r.Children[1:]
	for _, v := range l.Children {
		if v != nil {
			v.Parent = l
		}
	}
}

func (b *Bnode) RotateRight(idx int) {
	ridx := idx + 1
	l := b.Children[idx]
	r := b.Children[ridx]

	r.InsertV(b.Val[idx], 0)
	r.InsertN(l.Children[len(l.Children)-1], 0)

	b.Val[idx] = l.Val[len(l.Val)-1]
	l.Val = l.Val[:len(l.Val)-1]
	l.Children = l.Children[:len(l.Children)-1]
	for _, v := range r.Children {
		if v != nil {
			v.Parent = r
		}
	}
}

func (b *Bnode) Rebalance(n int) {
	if b == nil {
		return
	}
	if b.Parent == nil {
		return
	}
	if len(b.Children) >= n {
		return
	}
	fmt.Println("Rebalance ", b.Val)

	var idx int
	parent := b.Parent
	if len(parent.Val) == 0 {
		return
	}
	for i, v := range parent.Children {
		if v == b {
			idx = i
		}
	}
	left := true
	right := true

	if idx == 0 {
		left = false
		if len(parent.Children[idx+1].Children) > n {
			// rotate left
			fmt.Println("rotate left")
			parent.RotateLeft(idx)
			return
		}
	}
	if idx == len(parent.Children)-1 {
		right = false
		if len(parent.Children[idx-1].Children) > n {
			// rotate  right
			fmt.Println("rotate right")
			parent.RotateRight(idx - 1)
			return
		}
	}
	{
		if left && len(parent.Children[idx-1].Children) > n {
			// rotate right
			fmt.Println("rotate right")
			parent.RotateRight(idx - 1)
			return
		} else if right && len(parent.Children[idx+1].Children) > n {
			// rotate left
			fmt.Println("rotate left")
			parent.RotateLeft(idx)
			return
		} else {
			if left {
				// merge with left
				fmt.Println("merge left")

				v := parent.Val[idx-1]
				l := parent.Children[idx-1]
				r := b

				l.Val = append(l.Val, v)
				l.Val = append(l.Val, r.Val...)
				l.Children = append(l.Children, r.Children...)

				parent.Children = append(parent.Children[:idx], parent.Children[idx+1:]...)
				parent.Val = append(parent.Val[:idx-1], parent.Val[idx:]...)

				for _, v := range l.Children {
					if v != nil {
						v.Parent = l
					}
				}

				if len(parent.Children) < n {
					parent.Rebalance(n)
				}

			} else if right {
				// merge with right
				fmt.Println("merge right")
				v := parent.Val[idx]
				r := parent.Children[idx+1]
				l := b

				l.Val = append(l.Val, v)
				l.Val = append(l.Val, r.Val...)
				l.Children = append(l.Children, r.Children...)

				parent.Children = append(parent.Children[:idx+1], parent.Children[idx+2:]...)
				parent.Val = append(parent.Val[:idx], parent.Val[idx+1:]...)

				for _, v := range l.Children {
					if v != nil {
						v.Parent = l
					}
				}

				if len(parent.Children) < n {
					parent.Rebalance(n)
				}

			} else {
				// merge with nothing, what the fuck
				fmt.Println("merge with nothing,what the fuck!")
			}
		}
	}
	_, _ = left, right
}
