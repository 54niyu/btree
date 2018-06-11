package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
)

type Int int

func (i Int) Less(than Key) bool {
	return i < than.(Int)
}

func Test_BTP(t *testing.T) {

	b := NewBplusTree(4)

	items := make([]Key, 0)

	for i := 0; i < 30; i++ {

		r := rand.Intn(1000)
		b.Set(Int(r), struct {
			I int
			N string
		}{I: i, N: "fasdf"})
		//b.Print()
		items = append(items, Int(r))

		name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/a_%v_tree_%v.dot", i, r)
		b.Dot(name)
		cmd := exec.Command("dot",
			"-Tjpg", name, "-o",
			fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/b_%v_png_%v.jpg", i, r))
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}

	}

	for i, r := range items {
		fmt.Println(r)
		fmt.Printf("Get %v %v\n", r, b.Get(r))
		fmt.Println("Del ", r)
		b.Del(r)
		name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/a_%v_tree_del_%v.dot", i, r)

		b.Dot(name)
		cmd := exec.Command("dot",
			"-Tjpg", name, "-o",
			fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/c_%v_del_png_%v.jpg", i, r))
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}

	}

	for i := 0; i < 30; i++ {
		r := items[29-i]
		b.Set(r, i)
		name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/d_%v_tree_%v.dot", i, r)
		b.Dot(name)
		cmd := exec.Command("dot",
			"-Tjpg", name, "-o",
			fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/d_%v_png_%v.jpg", i, r))
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(items)
}
