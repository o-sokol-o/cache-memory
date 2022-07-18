package main

import (
	"fmt"
	"math"
	"time"

	cachememory "github.com/o-sokol-o/cache-memory"
)

func main() {
	fmt.Printf("\nNew cache = Life time cache value 3 sec\n")
	cache := cachememory.New(3)

	cache.Set("userId", 42)
	cache.Set("Pi", math.Pi)

	fmt.Printf("Get:  ")
	if userId, err := cache.Get("userId"); err == nil {
		fmt.Printf("Key = %s   Value = %v\n", "userId", userId)
	} else {
		fmt.Println(err.Error())
	}

	cache.Delete("userId")

	fmt.Printf("Get:  ")
	if userId, err := cache.Get("userId"); err == nil {
		fmt.Printf("Key = %s   Value = %v\n", "userId", userId)
	} else {
		fmt.Println(err.Error())
	}

	// - - - - - -

	fmt.Printf("\nLife time cache:\n")

	for i := 0; i < 6; i++ {

		fmt.Printf("%d sec --- Get:  ", i)
		if pi, err := cache.Get("Pi"); err == nil {
			fmt.Printf("Key = %s   Value = %v\n", "Pi", pi)
		} else {
			fmt.Println(err.Error())
		}

		time.Sleep(1000000000)
	}

	cache.Free()
}
