// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - go-jsschema https://github.com/lestrrat/go-jsschema
//     Copyright (c) 2016 lestrrat

package cui

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Arguments struct {
	OutputFileName string
	PackageName    string
	TypeName       string
	ShowVersion    bool
}

// Usage show help message.
func Usage() {
	fmt.Fprint(os.Stderr, "SYNOPSIS\n")
	fmt.Fprint(os.Stderr, "  schemarshal [options] [<json_schema_file>]\n")
	fmt.Fprint(os.Stderr, "OPTIONS\n")
	fmt.Fprint(os.Stderr, "  -h, -help\n")
	fmt.Fprint(os.Stderr, "           Show this help message.\n")
	fmt.Fprint(os.Stderr, "  -o <file>, -output <file>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("o").Usage)
	fmt.Fprint(os.Stderr, "  -p <package>, -package <package>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("p").Usage)
	fmt.Fprint(os.Stderr, "  -t <package>, -type <package>\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("t").Usage)
	fmt.Fprint(os.Stderr, "  -v, -version\n")
	fmt.Fprintf(os.Stderr, "           %s\n", flag.Lookup("v").Usage)
}

// ParseArguments parse command-line arguments
func ParseArguments() *Arguments {
	args := new(Arguments)
	flag.Usage = Usage
	flag.StringVar(&args.OutputFileName, "o", "", "Write output to file instead of stdout.")
	flag.StringVar(&args.OutputFileName, "output", "", "Write output to file instead of stdout.")
	flag.StringVar(&args.PackageName, "p", "main", "Package name for output. (default `main`)")
	flag.StringVar(&args.PackageName, "package", "main", "Package name for output. (default `main`)")
	flag.StringVar(&args.TypeName, "t", "", "Set default Type name.")
	flag.StringVar(&args.TypeName, "type", "", "Set default Type name.")
	flag.BoolVar(&args.ShowVersion, "v", false, "Show version.")
	flag.BoolVar(&args.ShowVersion, "version", false, "Show version.")
	flag.Parse()
	return args
}

// FileName return file-name without ext
func FileName(file *os.File) string {
	name := path.Base(file.Name())
	ext := path.Ext(name)
	return strings.TrimRight(name, ext)
}