package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "embed"
	
	datarecord "github.com/nikolainp/toGraph/datarecord"
	state "github.com/nikolainp/toGraph/statecontext"
)


//go:embed htmlTemplate.gohtml
var graphTemplate string


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

	config := state.Config()
	for _, fileName := range config.InputFiles {
		inputFile := getInputStream(fileName)
		outputFile, outName := getOutputStream(fileName, config)

		log.Printf("file being processed: %s", fileName)
		processFile(inputFile, outputFile, config, outName)
	}
}

func getInputStream(fileName string) io.Reader {
	if fileName == "" {
		return os.Stdin
	}

	inputFile, err := os.Open(fileName)
	checkErr(err)

	return inputFile
}

func getOutputStream(fileName string, config state.Configuration) (io.Writer, string) {
	var outputName string
	var outputPath string

	if len(config.OutputFile) == 0 {
		outputPath, outputName = getOutFileName(fileName)
	} else {
		_, outputName = getOutFileName(config.OutputFile)
		outputPath = config.OutputFile
	}
	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	checkErr(err)

	return outputFile, outputName
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
	reader.WithColumnNames(config.ColumnNames)
	reader.WithColumnNamesInFirstRow(config.IsColumnNamesInFirstRow)

	dataGraph, err := template.New("dataGraph").Parse(graphTemplate)
	checkErr(err)
	for i := 0; scanner.Scan(); i++ {
		dataString := scanner.Text()
		reader.ReadDataRecord(dataString)
	}

	data := struct {
		Title    string
		Columns  []datarecord.ColumnStatistic
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
