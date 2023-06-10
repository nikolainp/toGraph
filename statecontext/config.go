package statecontext

import (
	"flag"
	"fmt"
	"os"
)

func (obj *Config) configure(args ...string) {
	fs, err := readCommandLineArguments(obj, args)
	obj.printErrorAndUsage(fs, err)
}

func (obj *Config) printErrorAndUsage(fs *flag.FlagSet, err error) {
	if err != nil {
		fmt.Fprintf(fs.Output(), "Error:%s:\n", err)
	}
	if obj.printUsage {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()

		os.Exit(1)
	}
}

func readCommandLineArguments(config *Config, args []string) (fs *flag.FlagSet, err error) {
	fs = flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&config.printUsage, "h", false, "print usage")
	fs.StringVar(&config.NewLineRegex, "l", "", "reqular expression to determine the first line of an entry")
	fs.StringVar(&config.LogOutputPath, "o", "", "log output file")

	if err = fs.Parse(args); err != nil {
		config.printUsage = true
		return
	}

	switch fs.NArg() {
	case 1:
		config.SearchLineRegex = fs.Arg(0)
	case 2:
		config.SearchLineRegex = fs.Arg(0)
		config.SearchPath = fs.Arg(1)
	default:
		config.printUsage = true
	}

	return
}
