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

// copied from golint https://github.com/golang/lint/blob/c5fb716d6688a859aae56d26d3e6070808df29f7/lint.go#L742-L781
//// Copyright (c) 2013 The Go Authors. All rights reserved.
////
//// Use of this source code is governed by a BSD-style
//// license that can be found in the LICENSE file or at
//// https://developers.google.com/open-source/licenses/bsd.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// copy end

// UpperCamelCase returns the string to upper camel case
func UpperCamelCase(str string) string {
	matches := notAlphaNum.Split(str, -1)
	result := ""
	for i, m := range matches {
		var varWord string
		if i == 0 {
			varWord = alphaNum.FindString(m)
		} else {
			varWord = m
		}
		if u := strings.ToUpper(varWord); commonInitialisms[u] {
			result += u
		} else {
			result += strings.Title(varWord)
		}
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
