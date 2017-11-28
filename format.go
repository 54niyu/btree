package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

// PrintBtree print btree
func PrintBtree(t *Btree) {
	t.Root.PrintNode()
}

func (t *Bnode) PrintNode() {
	fmt.Println("--------Start------------")
	defer fmt.Println("--------End------------")
	lrtr := &Bnode{}
	sep := &Bnode{}
	queue := make([]*Bnode, 0)
	queue = append(queue, t)
	queue = append(queue, lrtr)

	for len(queue) != 0 {

		for i := 0; i < len(queue); i++ {
			if queue[i] == nil {
				//fmt.Print(" nil ")
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
				fmt.Print(queue[i].Val)
				queue = append(queue, queue[i].Children...)
				queue = append(queue, sep)
			}
		}
	}
}

// BeautifulPrint use gg to draw btree by own , very ugly
func (t *Bnode) BeautifulPrint(name string) {

	queue := make([][]*Bnode, 0)
	queue = append(queue, []*Bnode{t})

	for len(queue[len(queue)-1]) != 0 {
		q := queue[len(queue)-1]
		p := make([]*Bnode, 0)
		for i := 0; i < len(q); i++ {
			if q[i] == nil {
				//fmt.Print(" nil ")
			} else {
				p = append(p, q[i].Children...)
			}
		}
		queue = append(queue, p)
	}

	queue = queue[:len(queue)-1]
	each := 30
	width := len(queue[len(queue)-1]) * each
	hight := len(queue) * each

	dc := gg.NewContext(width+100, hight+100)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	for y, v := range queue {
		for x, v2 := range v {
			i := float64((y + 1) * each)
			j := float64((width) / (len(v) + 1) * (x + 1))
			if v2 != nil {
				dc.SetColor(color.Black)
				dc.DrawStringAnchored(fmt.Sprintf("%v", v2.Val), j, i, 0.5, 0.5)
				v2.x = j
				v2.y = i
				if v2.Parent != nil {
					dc.SetLineWidth(1)
					dc.DrawLine(float64(j), float64(i), v2.Parent.x, v2.Parent.y)
					dc.Stroke()
				}
			} else {
				dc.SetColor(color.Black)
				dc.DrawStringAnchored("[]", j, i, 0.5, 0.5)
			}
		}
	}
	dc.SavePNG(name + ".png")
}

// PrintTerminal printf btree in terminal
func (t *Btree) PrintTerminal() {
	t.Root.PrintTerminal(0)
}

func (t *Bnode) PrintTerminal(idx int) {
	for i, v := range t.Val {
		for j := 0; j < (idx-1)*5; j++ {
			fmt.Print(" ")
		}
		fmt.Print("|")
		for j := 0; j < 5; j++ {
			fmt.Print("-")
		}
		fmt.Println(v)
		if t.Children[i] != nil {
			t.Children[i].PrintTerminal(idx + 1)
		}
	}
	if t.Children[len(t.Children)-1] != nil {
		t.Children[len(t.Children)-1].PrintTerminal(idx + 1)
	}
}

// Dot generate dot file to draw btree
func (t *Bnode) Dot() {

	tpl := " digraph graphname {\n"

	lrtr := &Bnode{}
	queue := make([]*Bnode, 0)
	queue = append(queue, t, lrtr)

	for len(queue) != 0 {

		for i := 0; i < len(queue); i++ {
			if queue[i] == nil {
			} else if queue[i] == lrtr {
				queue = queue[i+1 : len(queue)]
				queue = append(queue, lrtr)
				if len(queue) == 1 {
					goto OVER
				}
				break
			} else {
				v := queue[i]
				var s []string
				for _, v := range v.Val {
					s = append(s, fmt.Sprintf("%v", v))
				}
				tpl += fmt.Sprintf("\t\"%p\" [label=\"%v\"];\n", v, strings.Join(s, "-"))
				if v.Parent != nil {
					tpl += fmt.Sprintf("\t\"%p\" ->  \"%p\" ;\n", v.Parent, v)
				}
				queue = append(queue, queue[i].Children...)
			}
		}
	}
OVER:
	tpl += "}"
	fmt.Println(tpl)
	ioutil.WriteFile("btree.dot", []byte(tpl), os.FileMode(0777))
}
