package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	run(os.Stdin, os.Stdout)
}

func run(sIn io.Reader, sOut io.Writer) error {
	fmt.Fprint(sOut, "test")
	return nil
}
