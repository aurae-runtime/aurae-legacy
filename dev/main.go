package main

import (
	"fmt"
	"github.com/kris-nova/aurae/pkg/core/memfs"
)

func main() {

	db := memfs.NewDatabase()
	db.Set("/beeps/boops/meeps/moops", "testvalue")
	result := db.Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		fmt.Println("FAILED")
		fmt.Println(result)
	}
}
