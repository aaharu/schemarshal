// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the 2-clause BSD license found in
// the LICENSE file in the root directory of this source tree.

package utils

import (
	"os"
	"testing"
)

func TestUpperCamelCase(t *testing.T) {
	actual := UpperCamelCase("1st address")
	expected := "StAddress"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = UpperCamelCase("a1st address")
	expected = "A1stAddress"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = UpperCamelCase("address 1st_url")
	expected = "Address1stURL"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = UpperCamelCase("quote\" slash/")
	expected = "QuoteSlash"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = UpperCamelCase("box1")
	expected = "Box1"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestCleanDescription(t *testing.T) {
	actual := CleanDescription(`改
行`)
	expected := "改 行"
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
