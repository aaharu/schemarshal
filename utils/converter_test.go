// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package utils

import (
	"os"
	"testing"
)

func TestUpperCamelCase(t *testing.T) {
	actual := UpperCamelCase("address")
	expected := "Address"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestFileName(t *testing.T) {
	file, _ := os.Open("./converter_test.go")
	defer file.Close()

	actual := FileName(file)
	expected := "converter_test"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestEnumTypeName(t *testing.T) {
	actual := EnumTypeName("address")
	expected := "AddressEnum"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
