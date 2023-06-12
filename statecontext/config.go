package statecontext

import (
	"flag"
	"fmt"
	"os"
)

var errEmptyArgumentList = fmt.Errorf("empty argument list")

type Configuration struct {
	InputFiles []string

	programName string
	printUsage  bool
}

func (obj *Configuration) configure(args []string) {
	fs, err := readCommandLineArguments(obj, args)
	obj.printErrorAndUsage(fs, err)
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

func readCommandLineArguments(config *Configuration, args []string) (fs *flag.FlagSet, err error) {
	fs = flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&config.printUsage, "h", false, "print usage")
	// fs.StringVar(&config.NewLineRegex, "l", "", "reqular expression to determine the first line of an entry")
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
