package main

import (
	"log"
	"sync"
)

type ttS struct {
	n int
}

func main() {
	var m sync.Map

	t := &ttS{n: 1}
	log.Println("before t=", t)
	m.Store(t.n, t)

	v, _ := m.Load(1)
	log.Println("load v=", v.(*ttS))
	t.n = 2

	log.Println("final v=", v)

	m.Range(func(k, v interface{}) bool {
		log.Println("range before v= ", v.(*ttS))
		t.n = 3
		log.Println("range after v=", v.(*ttS))
		return true
	})
}
