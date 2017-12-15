// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the 2-clause BSD license found in
// the LICENSE file in the root directory of this source tree.

package cui

import (
	"flag"
	"fmt"
	"os"
)

// Arguments are CUI options
type Arguments struct {
	InputFileName  string
	OutputFileName string
	PackageName    string
	TypeName       string
	ShowVersion    bool
	NoComment      bool
}

// Usage show help message.
func Usage() {
	fmt.Fprint(os.Stderr, "SYNOPSIS\n")
	fmt.Fprint(os.Stderr, "  schemarshal [options] [<json_schema_file>]\n")
	fmt.Fprint(os.Stderr, "OPTIONS\n")
	fmt.Fprint(os.Stderr, "  -h, -help\n")
	fmt.Fprint(os.Stderr, "           Show this help message.\n")
	fmt.Fprint(os.Stderr, "  -f <file>, -file <file>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("f").Usage)
	fmt.Fprint(os.Stderr, "  -o <file>, -output <file>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("o").Usage)
	fmt.Fprint(os.Stderr, "  -p <package>, -package <package>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("p").Usage)
	fmt.Fprint(os.Stderr, "  -t <package>, -type <package>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("t").Usage)
	fmt.Fprint(os.Stderr, "  -v, -version\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("v").Usage)
	fmt.Fprint(os.Stderr, "  -nc, -nocomment\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("nc").Usage)
}

// ParseArguments parse command-line arguments
func ParseArguments() *Arguments {
	args := &Arguments{}
	flag.Usage = Usage
	flag.StringVar(&args.InputFileName, "f", "", "Input file name.")
	flag.StringVar(&args.InputFileName, "file", "", "Input file name.")
	flag.StringVar(&args.OutputFileName, "o", "", "Write output to file instead of stdout.")
	flag.StringVar(&args.OutputFileName, "output", "", "Write output to file instead of stdout.")
	flag.StringVar(&args.PackageName, "p", "main", "Package name for output. (default `main`)")
	flag.StringVar(&args.PackageName, "package", "main", "Package name for output. (default `main`)")
	flag.StringVar(&args.TypeName, "t", "", "Set default Type name.")
	flag.StringVar(&args.TypeName, "type", "", "Set default Type name.")
	flag.BoolVar(&args.ShowVersion, "v", false, "Show version.")
	flag.BoolVar(&args.ShowVersion, "version", false, "Show version.")
	flag.BoolVar(&args.NoComment, "nc", false, "Do not output comments.")
	flag.BoolVar(&args.NoComment, "nocomment", false, "Do not output comments.")
	flag.Parse()
	return args
}
