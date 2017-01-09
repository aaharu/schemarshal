// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package utils

import (
	"os"
	"path"
	"strings"
)

// Ucfirst Upper case first character
func Ucfirst(str string) string {
	return strings.Replace(strings.Title(str), " ", "", -1)
}

// FileName return file-name without ext
func FileName(file *os.File) string {
	name := path.Base(file.Name())
	ext := path.Ext(name)
	return strings.TrimRight(name, ext)
}
