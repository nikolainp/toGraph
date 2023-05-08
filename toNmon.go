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

const header string = `{
"hostname": "dummy.rtp.raleigh.ibm.com",
"whenPattern": "MM/dd/yyyy HH:mm:ss",
"metadata": {
	"mdata1": "value 1",
	"mdata2": "value 2",
	"mdata3": "completely arbitrary data goes here"
},
"types": [
	{
	"id": "type1",
	"name": "type 1",
	"fields": [ "field11", "field12", "field13" ]
	}
],
"data": [
`

func run(sIn io.Reader, sOut io.Writer) error {
	scanner := bufio.NewScanner(sIn)
	writer := bufio.NewWriter(sOut)

	writer.WriteString(header)
	for i := 0; scanner.Scan(); i++ {
		data := scanner.Text()
		record := datarecord.GetDataRecord(data)

		if i == 0 {
			writer.WriteString(record.StringWithIndention(1))
		} else {
			writer.WriteString(fmt.Sprintf(",\n%s", record.StringWithIndention(1)))
		}
	}
	writer.WriteString("\n]\n}\n")
	writer.Flush()

	return nil
}
