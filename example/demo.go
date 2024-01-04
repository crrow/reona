package main

import (
	"fmt"
	"sync"

	"github.com/crrow/reona/linkedlist"
)

func main() {
	mem := linkedlist.NewMap[string, int](linkedlist.WithCapacity[string, int](10))
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		for {
			r, ok := mem.Get("hello")
			if !ok {
				fmt.Println("cannot get hello1")
				continue
			}
			fmt.Println(*r)
			return
		}
	}()

	_, ok := mem.Get("hello")
	fmt.Println(ok)
	_, ok = mem.Get("hello2")
	fmt.Println(ok)

	go func() {
		defer wg.Done()
		mem.Insert("hello", 1)
	}()

	go func() {
		defer wg.Done()
		mem.Insert("hello2", 2)
	}()

	go func() {
		defer wg.Done()
		for {
			r, ok := mem.Get("hello2")
			if !ok {
				fmt.Println("cannot get hello2")
				continue
			}
			fmt.Println(*r)
			return
		}
	}()

	wg.Wait()
}
