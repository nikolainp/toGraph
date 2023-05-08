package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	datarecord "github.com/nikolainp/toNmon/datarecord"
)

func main() {
	run(os.Stdin, os.Stdout)
}

func run(sIn io.Reader, sOut io.Writer) error {
	fmt.Fprint(sOut, "test")

	scanner := bufio.NewScanner(sIn)
	for scanner.Scan() {
		data := scanner.Text()
		fmt.Fprint(sOut, data)
		record := datarecord.GetDataRecord(data)
		fmt.Fprint(sOut, record)
	}
	return nil
}
