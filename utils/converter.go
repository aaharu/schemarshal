// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package utils

import (
	"os"
	"path"
	"strings"
)

// UpperCamelCase returns the string to upper camel case
func UpperCamelCase(str string) string {
	str = strings.Replace(str, "-", " ", -1)
	str = strings.Replace(str, "_", " ", -1)
	return strings.Replace(strings.Title(str), " ", "", -1)
}

// FileName returns file-name without ext
func FileName(file *os.File) string {
	name := path.Base(file.Name())
	ext := path.Ext(name)
	return strings.TrimRight(name, ext)
}

// EnumTypeName returns Go Type literals
func EnumTypeName(str string) string {
	return UpperCamelCase(str) + "Enum"
}
