package main

import (
	"log"
)

type ttS struct {
	n int
}

func main() {
	mapTest := make(map[int]*ttS)

	t := &ttS{n: 1}
	log.Println("before t=", t)
	mapTest[t.n] = t

	v, _ := mapTest[1]
	log.Println("load v=", v)
	f := *v
	v.n = 2

	log.Println("final v=", v, ", f=", f)

}
