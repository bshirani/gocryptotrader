package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("ok\n")
	panic("a problem")

	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}
}
