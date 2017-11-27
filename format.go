package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

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

func (t *Btree) DrawPic(path string) {

	dc := gg.NewContext(10000, 1000)
	dc.SetRGB(1, 1, 1)
	dc.Fill()
	t.Root.DrawPic(dc, 5000, 0, 1, 5000)
	dc.SavePNG(path)

}

func (t *Bnode) DrawPic(dc *gg.Context, x, y float64, idx int, with float64) {

	sf := float64(10-idx) / 10.0

	dc.DrawCircle(x, y+20, 20*sf)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SetRGB(1, 0, 0)
	dc.DrawString(fmt.Sprintf("%v", t.Val), x, y+20)

	step := float64(with*2) / float64(len(t.Children)+1)

	for i, v := range t.Children {
		if v != nil {
			nx := x - with + float64(i+1)*step
			ny := y + 100
			dc.SetRGB(0, 0, 0)
			dc.SetLineWidth(2)
			dc.DrawLine(x, y+20, nx, ny)
			dc.Stroke()
		}
	}

	for i, v := range t.Children {
		if v != nil {
			v.DrawPic(dc, x-with+float64(i+1)*step, y+100, idx+1, step/2)
		}
	}
}

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

func (t *Bnode) Dot() {

	lab := 0
	lrtr := &Bnode{}
	queue := make([]*Bnode, 0)
	queue = append(queue, t)
	queue = append(queue, lrtr)
	record := make([]*Bnode, 0)

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
				fmt.Print(queue[i].Val)
				queue[i].i = lab
				lab++
				queue = append(queue, queue[i].Children...)
				record = append(record, queue[i])
			}
		}
	}
OVER:

	fmt.Println("")
	for _, v := range record {
		for _, v2 := range v.Children {
			if v2 != nil {
				var s []string
				for _, v := range v.Val {
					s = append(s, fmt.Sprintf("%v", v))
				}
				fmt.Printf("%v [label=\"%v\"];\n", v.i, strings.Join(s, "-"))
				fmt.Printf("%v -> %v ;\n", v.i, v2.i)
			}
		}
	}

	fmt.Println("---", len(record))

}
