package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	cachememory "github.com/o-sokol-o/cache-memory"
)

func main() {
	fmt.Printf("\nNew cache = Life time cache value 3 sec\n")
	cache := cachememory.New(cachememory.ResolutionSeconds)

	cache.Set("userId", 42, 0)
	cache.Set("Pi", math.Pi, 0)

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

	// - - - - - - Example lifetime - - - - - - -

	fmt.Printf("\nLife time cache:\n")

	for i := 0; i < 6; i++ {

		fmt.Printf("%d sec --- Get:  ", i)
		if pi, err := cache.Get("Pi"); err == nil {
			fmt.Printf("Key = %s   Value = %v\n", "Pi", pi)
		} else {
			fmt.Println(err.Error())
		}

		time.Sleep(time.Second)
	}

	// - - - - - - Example race condition: go run -race example/main.go  - - - - - - -

	fmt.Printf("\n8000 goroutines write cache\n")
	for i := 0; i < 8000; i++ {
		go func() {
			cache.Set("userId"+strconv.Itoa(rand.Intn(10000000)), 42, 7)
		}()
	}

	fmt.Printf("Sleep 10 Second\n")
	time.Sleep(10 * time.Second)

	fmt.Printf("cache.Free\n")
	cache.Free()
}
