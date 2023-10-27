package main

import (
	"bufio"
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

// TODO: column names
// TODO: graph name + output file name
// TODO: statistic by columns

func main() {
	state.InitState()
	state.Configure(os.Args)

	run()
}

func run() {

	for _, fileName := range state.Config().InputFiles {
		inputFile, err := os.Open(fileName)
		checkErr(err)

		outputFile, err := os.OpenFile(fileName+".html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
		checkErr(err)

		log.Printf("file being processed: %s", fileName)
		processFile(inputFile, outputFile, state.Config())
	}
}

func processFile(sIn io.Reader, sOut io.Writer, config state.Configuration) error {

	scanner := bufio.NewScanner(sIn)
	reader := datarecord.GetDataReader()
	reader.WithDateFormat(config.DateFormat)
	reader.WithDateColumn(config.DateColumn)
	reader.WithPivotColumn(config.PivotColumn)
	reader.WithDelimiter(config.Delimiter)

	dataGraph, err := template.New("dataGraph").Parse(graphTemplate)
	checkErr(err)
	for i := 0; scanner.Scan(); i++ {
		dataString := scanner.Text()
		reader.ReadDataRecord(dataString)
	}

	data := struct {
		Title    string
		Columns  []string
		DataRows []string
	}{
		Title:    "My page",
		Columns:  reader.GetColumns(),
		DataRows: reader.GetDataRows(),
	}

	err = dataGraph.Execute(sOut, data)
	checkErr(err)

	return nil
}
