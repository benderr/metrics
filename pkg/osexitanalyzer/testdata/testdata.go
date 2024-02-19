package testdata

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("%s", 1)

	os.Exit(1) // want "preventing call os.Exit in main function"

	defer func() {
		os.Exit(1) // want "preventing call os.Exit in main function"
	}()
}

func main2() {
	os.Exit(2) // ignored
}
