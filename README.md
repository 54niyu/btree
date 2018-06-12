### golang btree data struct implementation

## btree 

## b+tree

```Golang
type Int int

func (i Int) Less(than Key) bool {
	return i < than.(Int)
}

func main(){

	b := NewBplusTree(4)

	items := make([]Key, 0)

	for i := 0; i < 30; i++ {
		r := rand.Intn(1000)
		b.Set(Int(r), struct {
			I int
			N string
		}{I: i, N: "fasdf"})
		items = append(items, Int(r))
	}

	b.Iter(func(key Key, val interface{}) error {
		fmt.Println("In interator ", key, val)
		if !key.Less(Int(258)) && !Int(258).Less(key) {
			return fmt.Errorf("Stop")
		}
		return nil
	})

	for i, r := range items {
		fmt.Printf("Get %v %v\n", r, b.Get(r))
		fmt.Println("Del ", r)
		b.Del(r)
	}
}

```
