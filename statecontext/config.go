package statecontext

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var errEmptyArgumentList = fmt.Errorf("empty argument list")

type Configuration struct {
	InputFiles []string
	DateFormat string

	programName string
	printUsage  bool
}

func (obj *Configuration) configure(args []string) {
	fs, err := readCommandLineArguments(obj, args)
	obj.printErrorAndUsage(fs, err)
	obj.DateFormat = covertDateFormat(obj.DateFormat)
}

func (obj *Configuration) printErrorAndUsage(fs *flag.FlagSet, err error) {
	if err != nil {
		fmt.Fprintf(fs.Output(), "Error:%s:\n", err)
	}
	if obj.printUsage {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", obj.programName)
		fs.PrintDefaults()

		os.Exit(1)
	}
}

func covertDateFormat(dateFormat string) string {
	// "YYYYMMDDHHMMSS", "20060102150405"
	dateFormat = strings.Replace(dateFormat, "YYYY", "2006", 1)
	dateFormat = strings.Replace(dateFormat, "YY", "06", 1)
	dateFormat = strings.Replace(dateFormat, "MM", "01", 1)
	dateFormat = strings.Replace(dateFormat, "DD", "02", 1)
	dateFormat = strings.Replace(dateFormat, "HH", "15", 1)
	dateFormat = strings.Replace(dateFormat, "MM", "04", 1)
	dateFormat = strings.Replace(dateFormat, "SS", "05", 1)

	return dateFormat
}

func readCommandLineArguments(config *Configuration, args []string) (fs *flag.FlagSet, err error) {
	fs = flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&config.printUsage, "h", false, "print usage")
	fs.StringVar(&config.DateFormat, "d", "YYYYMMDDHHMMSS", "time field format (YYYY-MM-DDTHH:MM:SS.mmmmmm), YYYYMMDDHHMMSS by default")
	// fs.StringVar(&config.LogOutputPath, "o", "", "log output file")

	if len(args) == 0 {
		return nil, errEmptyArgumentList
	}

	config.programName = args[0]
	if err = fs.Parse(args[1:]); err != nil {
		config.printUsage = true
		return
	}

	switch {
	case 0 < fs.NArg():
		config.InputFiles = fs.Args()
	default:
		config.printUsage = true
	}

	return
}
