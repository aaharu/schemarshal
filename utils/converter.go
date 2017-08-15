// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package utils

import (
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	alphaNum    = regexp.MustCompile(`[a-zA-Z]+[0-9a-zA-Z]*`)
	notAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// UpperCamelCase returns the string to upper camel case
func UpperCamelCase(str string) string {
	matches := notAlphaNum.Split(str, -1)
	result := ""
	for i, m := range matches {
		if i == 0 {
			result += strings.Title(alphaNum.FindString(m))
			continue
		}
		result += strings.Title(m)
	}
	return result
}

// CleanDescription remove \n, \r and \t
func CleanDescription(desc string) string {
	desc = strings.TrimSpace(desc)
	desc = strings.Replace(desc, "\n", " ", -1)
	desc = strings.Replace(desc, "\r", " ", -1)
	desc = strings.Replace(desc, "\t", " ", -1)
	return desc
}

// FileName returns file-name without ext
func FileName(file *os.File) string {
	name := path.Base(file.Name())
	ext := path.Ext(name)
	return name[0 : len(name)-len(ext)]
}
