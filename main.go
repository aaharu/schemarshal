// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

// This source code use following software(s):
//   - go-jsschema https://github.com/lestrrat/go-jsschema
//     Copyright (c) 2016 lestrrat

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

	schema "github.com/lestrrat/go-jsschema"

	"github.com/aaharu/schemarshal/cui"
	"github.com/aaharu/schemarshal/jsonschema"
	"github.com/aaharu/schemarshal/version"
)

func main() {
	args := cui.ParseArguments()

	if args.ShowVersion {
		fmt.Printf("schemarshal %s\n", version.Version)
		os.Exit(0)
	}

	var input io.Reader
	typeName := args.TypeName
	if terminal.IsTerminal(syscall.Stdin) {
		if len(flag.Args()) < 1 {
			cui.Usage()
			os.Exit(1)
		}

		inputFile, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open schema: %s\n", err)
			os.Exit(1)
		}
		defer inputFile.Close()

		if typeName == "" {
			typeName = cui.FileName(inputFile)
		}
		input = inputFile
	} else {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
			os.Exit(1)
		}

		if typeName == "" {
			typeName = "T"
		}
		input = strings.NewReader(string(stdin))
	}

	jsschema, err := schema.Read(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schema: %s\n", err)
		os.Exit(1)
	}

	js := jsonschema.New(jsschema)
	js.SetCommand(strings.Trim(fmt.Sprintf("%v", os.Args), "[]"))
	js.SetPackageName(args.PackageName)
	output, err := js.Typedef(typeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate: %s\n", err)
		os.Exit(1)
	}

	if args.OutputFileName == "" {
		fmt.Printf("%s\n", output)
	} else {
		outputFile, err := os.Create(args.OutputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file: %s\n", err)
			os.Exit(1)
		}
		defer outputFile.Close()

		outputFile.Write([]byte(output))
	}
}
