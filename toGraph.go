package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

		outFilePath, outFileName := getOutFileName(fileName)

		outputFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
		checkErr(err)

		log.Printf("file being processed: %s", fileName)
		processFile(inputFile, outputFile, state.Config(), outFileName)
	}
}

func getOutFileName(fileName string) (outFilePath string, outFileName string) {
	outFilePath = filepath.Dir(fileName)
	outFileName = filepath.Base(fileName)
	outFileName = strings.TrimSuffix(outFileName, filepath.Ext(outFileName))

	return filepath.Join(outFilePath, outFileName+".html"), outFileName
}

func processFile(sIn io.Reader, sOut io.Writer, config state.Configuration, title string) error {

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
		Title:    title,
		Columns:  reader.GetColumns(),
		DataRows: reader.GetDataRows(),
	}

	err = dataGraph.Execute(sOut, data)
	checkErr(err)

	return nil
}
