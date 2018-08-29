package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

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
