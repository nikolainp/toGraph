package main

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"

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
	scanner := bufio.NewScanner(sIn)
	//writer := bufio.NewWriter(sOut)

	dataGraph, err := template.New("dataGraph").Parse(graphTemplate)
	checkErr(err)
	for i := 0; scanner.Scan(); i++ {
		data := scanner.Text()
		record := datarecord.GetDataRecord(data)
		_ = record
		//if i == 0 {
		//	writer.WriteString(record.StringWithIndention(1))
		//} else {
		//	writer.WriteString(fmt.Sprintf(",\n%s", record.StringWithIndention(1)))
		//}
	}

	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}

	err = dataGraph.Execute(os.Stdout, data)
	checkErr(err)

	return nil
}
