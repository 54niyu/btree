package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"
)

type Int int

func (i Int) Less(than Key) bool {
	return i < than.(Int)
}

func Test_RBTree(t *testing.T) {

	var saves []int64
	var num int = 100

	b := NewRBTree()
	rand.Seed(int64(time.Now().Nanosecond()))
	for idx := 0; idx < num; idx++ {
		r := rand.Intn(1000)
		saves = append(saves, int64(r))
		b.Set(Int(r), idx)
		b.I()
		fmt.Println("")
		b.Hight()
		fmt.Println("")
		name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/rbtree_%v_tree_%v.dot", idx, r)
		b.Dot(name)
		cmd := exec.Command("dot",
			"-Tjpg", name, "-o",
			fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/rbtree_%v_png_%v.jpg", idx, r))
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	for i, v := range saves {
		b.Del(Int(v))
		b.I()
		fmt.Println("")
		b.Hight()
		fmt.Println("")
		name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/rbtree_del_%v_tree_%v.dot", i, v)
		b.Dot(name)
		cmd := exec.Command("dot",
			"-Tjpg", name, "-o",
			fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/rbtree_del_%v_png_%v.jpg", i, v))
		cmd.Env = os.Environ()
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}

	}
}
