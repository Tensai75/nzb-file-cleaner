package main

import (
	"fmt"
	"os"

	parser "github.com/alexflint/go-arg"
)

// arguments structure
type Args struct {
	NZBFile              string `arg:"positional,required" help:"path to the NZB file to be cleanded or a folder containing NZB files (required)"`
	DestPath             string `arg:"positional" help:"destination path where the new NZB file(s) should be saved (optional)"`
	AddPwToMeta          bool   `arg:"--apm" help:"add password from filename ({{password}}) to NZB file metadata"`
	AddPwToFilename      bool   `arg:"--apf" help:"add password from NZB file metadata to filename ({{password}})"`
	AddTitleToMeta       bool   `arg:"--atm" help:"add the filename to NZB file metadata as title"`
	UseTitleForFilename  bool   `arg:"--utf" help:"use the title in the NZB file metadata as the filename for the NZB file"`
	RemovePwFromMeta     bool   `arg:"--rpm" help:"remove password from the NZB file metadata"`
	RemovePwFromFilename bool   `arg:"--rpf" help:"remove password from the filename ({{password}})"`
	RemoveTitleFromMeta  bool   `arg:"--rtm" help:"remove the title from the NZB file metadata"`
	Verbose              bool   `arg:"-v,--verbose" help:"enable verbose output"`
}

// version information
func (Args) Version() string {
	return fmt.Sprintf("%v %v", appName, appVersion)
}

// additional description
func (Args) Epilogue() string {
	return "For more information visit github.com/Tensai75/nzb-file-cleaner\n"
}

// global arguments variable
var args struct {
	Args
}

// parser variable
func parseArguments() {
	var argParser *parser.Parser

	parserConfig := parser.Config{
		IgnoreEnv: true,
	}

	// parse flags
	argParser, _ = parser.NewParser(parserConfig, &args)
	if err := parser.Parse(&args); err != nil {
		if err.Error() == "help requested by user" {
			argParser.WriteHelp(os.Stdout)
			os.Exit(0)
		} else if err.Error() == "version requested by user" {
			fmt.Println(args.Version())
			os.Exit(0)
		}
		argParser.WriteHelp(os.Stdout)
		exit(err)
	}

	checkArguments(argParser)

}

func checkArguments(argParser *parser.Parser) {
	// Exit with an error if no arguments are passed
	if len(os.Args) <= 1 {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("no arguments provided"))
	}

	// Exit with an error if no NZBFile argument is provided
	if args.NZBFile == "" {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("no NZBFile argument was given"))
	}

	// Exit with an error if only the NZBFile argument is provided
	if len(os.Args) == 2 && args.NZBFile != "" {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("no operation flags provided, only NZBFile argument was given"))
	}

	if args.AddPwToMeta && args.RemovePwFromMeta {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("cannot add and remove password to/from metadata at the same time"))
	}
	if args.AddPwToFilename && args.RemovePwFromFilename {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("cannot add and remove password to/from filename at the same time"))
	}
	if args.AddTitleToMeta && args.RemoveTitleFromMeta {
		argParser.WriteHelp(os.Stdout)
		exit(fmt.Errorf("cannot add and remove title to/from metadata at the same time"))
	}

}
