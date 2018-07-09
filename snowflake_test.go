package main

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan int64)
	count := 10000
	// 并发 count 个 goroutine 进行 snowflake ID 生成
	for i := 0; i < count; i++ {
		go func() {
			id := node.Generate()
			t.Log(id)
			ch <- id
		}()
	}

	defer close(ch)

	m := make(map[int64]int)
	for i := 0; i < count; i++  {
		id := <- ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m[id]
		if ok {
			t.Log("ID is not unique!\n")
			return
		}
		m[id] = i
	}
	t.Log("All ", count, " snowflake ID generate successed!\n")
}
