package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	datarecord "github.com/nikolainp/toGraph/datarecord"
	state "github.com/nikolainp/toGraph/statecontext"
)

var checkErr = func(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	state.InitState()
	state.Configure(os.Args)

	run()
}

func run() {

	for _, fileName := range state.Config().InputFiles {
		inputFile, err := os.Open(fileName)
		checkErr(err)

		outputFile, err := os.OpenFile(fileName+".html", os.O_CREATE|os.O_WRONLY, 0660)
		checkErr(err)

		processFile(inputFile, outputFile)
	}
}

func processFile(sIn io.Reader, sOut io.Writer) error {
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
