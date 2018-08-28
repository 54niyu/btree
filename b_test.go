package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func Test_BplusTree(t *testing.T) {

	/* tree := NewbPlusTree(3)*/
	//for i := 0; i < 50; i++ {
	//r := rand.Intn(1000)
	//name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/a_%v_tree_%v.dot", i, r)
	//fmt.Println(tree.Insert(r, r))
	//tree.Root.Print()
	//tree.Root.Dot(name)
	//cmd := exec.Command("dot",
	//"-Tjpg", name, "-o",
	//fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/b_%v_png_%v.jpg", i, r))
	//cmd.Env = os.Environ()
	//fmt.Println(cmd.Args)
	//err := cmd.Run()
	//if err != nil {
	//fmt.Println(err)
	//}
	//}

}

func Test_Btree(t *testing.T) {
	tree := NewBtree(20)
	var record []int
	var records []string

	rand.Seed(int64(time.Now().Nanosecond()))
	for i := 0; i < 50; i++ {
		rd := rand.Intn(1000)
		record = append(record, rd)
		fmt.Println("Insert ", rd)
		tree.Insert(rd)
		records = append(records, fmt.Sprintf("%v", rd))
	}
	for i := 0; i < 150; i++ {
		tree.Insert(i)
	}

	tree.Root.Dot()
	tree.Root.PrintNode()
	s := strings.Join(records, ",")
	fmt.Printf("(%v)\n", s)
	tree.Root.BeautifulPrint("test")

	for i := 0; i < 50; i++ {
		tree.Delete(record[i])
		tree.Root.PrintNode()
		//tree.Root.BeautifulPrint(fmt.Sprintf("%v_%v", i, record[i]))
	}
	for i := 49; i >= 0; i-- {
		tree.Delete(i)
		tree.Root.PrintNode()
		//tree.Root.BeautifulPrint(fmt.Sprintf("%v_%v", 100-i, i))
	}
	tree.Root.BeautifulPrint("finish")
}

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
		b.root.I()
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
	/*for i, v := range saves {*/
	//b.Del(Int(v))
	//b.root.I()
	//fmt.Println("")
	//name := fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/dot/rbtree_del_%v_tree_%v.dot", i, v)
	//b.Dot(name)
	//cmd := exec.Command("dot",
	//"-Tjpg", name, "-o",
	//fmt.Sprintf("/Users/Bing/Golang/src/learnGo/btree/img/rbtree_del_%v_png_%v.jpg", i, v))
	//cmd.Env = os.Environ()
	//err := cmd.Run()
	//if err != nil {
	//fmt.Println(err)
	//}

	/*}*/
}
