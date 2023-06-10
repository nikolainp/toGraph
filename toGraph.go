package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	datarecord "github.com/nikolainp/toGraph/datarecord"
)

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	run(os.Stdin, os.Stdout)
}

func run(sIn io.Reader, sOut io.Writer) error {
	data := struct {
		Title    string
		Columns  []string
		DataRows []string
	}{
		Title:    "My page",
		Columns:  []string{},
		DataRows: []string{},
	}

	scanner := bufio.NewScanner(sIn)

	dataGraph, err := template.New("dataGraph").Parse(graphTemplate)
	checkErr(err)
	for i := 0; scanner.Scan(); i++ {
		dataString := scanner.Text()
		record := datarecord.GetDataRecord(dataString)

		if i == 0 {
			for j := 0; j < record.Columns(); j++ {
				data.Columns = append(data.Columns, fmt.Sprintf("Column %d", j+1))
			}
		}

		data.DataRows = append(data.DataRows, record.String())
	}

	err = dataGraph.Execute(sOut, data)
	checkErr(err)

	return nil
}
