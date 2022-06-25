package main

import (
	"fmt"
	"github.com/kris-nova/aurae/pkg/core/mem"
)

func main() {

	db := mem.NewDatabase()
	db.Set("/beeps/boops/meeps/moops", "testvalue")
	result := db.Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		fmt.Println("FAILED")
		fmt.Println(result)
	}
}
