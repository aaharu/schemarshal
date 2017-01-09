// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package codegen

import (
	"reflect"
	"testing"
)

func TestJSONTagOmitEmpty(t *testing.T) {
	tag := &JSONTag{
		OmitEmpty: true,
	}
	actual := tag.Generate()
	expected := []byte("`json:\",omitempty\"`")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", string(actual), string(expected))
	}
}

func TestJSONTag(t *testing.T) {
	tag := &JSONTag{
		Name: "key",
	}
	actual := tag.Generate()
	expected := []byte("`json:\"key\"`")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
