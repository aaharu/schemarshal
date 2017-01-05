// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package schemarshal

import (
	"fmt"
	"strings"
)

// MarshalTag return tag for MarshalJSON
func MarshalTag(name string, required bool) string {
	if required {
		return "`json:\"" + name + "\"`"
	}
	return "`json:\"" + name + ",omitempty\"`"
}

// Ucfirst Upper case first character
func Ucfirst(str string) string {
	return strings.Replace(strings.Title(str), " ", "", -1)
}

// GeneratedByComment generate auto-generated comments for Go source code
func GeneratedByComment(command string) string {
	return fmt.Sprintf("// generated by schemarshal %s `%s`\n", Version, command)
}
