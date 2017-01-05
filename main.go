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
	"os"
	"path"
	"strings"

	schema "github.com/lestrrat/go-jsschema"

	jsonschema "github.com/aaharu/schemarshal/jsonschema"
	version "github.com/aaharu/schemarshal/version"
)

func main() {
	var (
		outputFileName string
		packageName    string
	)
	flag.StringVar(&outputFileName, "o", "", "Write output to file instead of stdout.")
	flag.StringVar(&outputFileName, "output", "", "Write output to file instead of stdout.")
	flag.StringVar(&packageName, "p", "main", "Package name for output.")
	flag.StringVar(&packageName, "package", "main", "Package name for output.")
	showVersion := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if len(os.Args) > 1 && *showVersion {
		fmt.Printf("schemarshal %s\n", version.Version)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		usage()
		os.Exit(1)
	}

	inputFile, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open schema: %s\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	jsschema, err := schema.Read(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schema: %s\n", err)
		os.Exit(1)
	}

	js := jsonschema.New(jsschema)
	js.SetCommand(strings.Trim(fmt.Sprintf("%v", os.Args), "[]"))
	js.SetPackageName(packageName)
	output, err := js.Typedef(fileName(inputFile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate: %s\n", err)
		os.Exit(1)
	}

	if outputFileName == "" {
		fmt.Printf("%s\n", output)
	} else {
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file: %s\n", err)
			os.Exit(1)
		}
		defer outputFile.Close()

		outputFile.Write([]byte(output))
	}
}

func usage() {
	fmt.Println("SYNOPSIS")
	fmt.Println("  schemarshal [options] <json_shcema_file>")
	fmt.Println("OPTIONS")
	fmt.Println("  -h, -help")
	fmt.Println("           Show help message.")
	fmt.Println("  -o <file>, -output <file>")
	fmt.Println("           Write output to <file> instead of stdout.")
	fmt.Println("  -p <package>, -package <package>")
	fmt.Println("           Package name for output. (default `main`)")
	fmt.Println("  -version")
	fmt.Println("           Show version.")
}

func fileName(file *os.File) string {
	name := path.Base(file.Name())
	ext := path.Ext(name)
	return strings.TrimRight(name, ext)
}
