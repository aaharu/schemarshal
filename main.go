// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - golang.org/x/crypto/ssh/terminal
//     Copyright 2011 The Go Authors.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/aaharu/schemarshal/codegen"
	"github.com/aaharu/schemarshal/cui"
	"github.com/aaharu/schemarshal/utils"
	"github.com/aaharu/schemarshal/version"
)

func main() {
	args := cui.ParseArguments()

	if args.ShowVersion {
		fmt.Println(version.String())
		os.Exit(0)
	}

	var input io.Reader
	typeName := args.TypeName
	if args.InputFileName != "" || terminal.IsTerminal(syscall.Stdin) {
		// input from file
		if args.InputFileName == "" && len(flag.Args()) < 1 {
			cui.Usage()
			os.Exit(1)
		}

		inputFileName := args.InputFileName
		if inputFileName == "" {
			inputFileName = flag.Args()[0]
		}
		inputFile, err := os.Open(inputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open schema: %s\n", err)
			os.Exit(1)
		}
		defer inputFile.Close()

		if typeName == "" {
			typeName = utils.FileName(inputFile)
		}
		input = inputFile
	} else {
		// input from pipe
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
			os.Exit(1)
		}

		if typeName == "" {
			// default type name
			typeName = "SchemarshalType"
		}
		input = strings.NewReader(string(stdin))
	}

	js, err := codegen.ReadSchema(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schema: %s\n", err)
		os.Exit(1)
	}

	codeGenerator := codegen.NewGenerator(args.PackageName, strings.Trim(fmt.Sprintf("%v", os.Args), "[]"))

	if js.GetTitle() != "" {
		typeName = js.GetTitle()
	}
	typeName = utils.UpperCamelCase(typeName)
	err = codeGenerator.AddSchema(typeName, js)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse: %s\n", err)
		os.Exit(1)
	}

	src, err := codeGenerator.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate: %s\n", err)
		os.Exit(1)
	}

	if args.OutputFileName == "" {
		fmt.Printf("%s\n", src)
	} else {
		outputFile, err := os.Create(args.OutputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file: %s\n", err)
			os.Exit(1)
		}
		outputFile.Write(src)
		defer outputFile.Close()
	}
}
