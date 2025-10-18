package main

import (
	"fmt"
	"os"

	"github.com/google/pprof/driver"
)

func main() {
	fmt.Println("Starting PProf Server")

	options := &driver.Options{}

	err := driver.PProf(options)
	if err != nil {
		fmt.Printf("error running PProf Server: %v\n", err)
		os.Exit(1)
	}
}
