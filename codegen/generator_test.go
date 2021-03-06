// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the 2-clause BSD license found in
// the LICENSE file in the root directory of this source tree.

package codegen

import (
	"os"
	"reflect"
	"testing"

	"github.com/aaharu/schemarshal/utils"
)

func TestJSONTagOmitEmpty(t *testing.T) {
	tag := &jsonTag{
		omitEmpty: true,
	}
	actual := tag.generate()
	expected := []byte("`json:\",omitempty\"`")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", string(actual), string(expected))
	}
}

func TestJSONTag(t *testing.T) {
	tag := &jsonTag{
		name: "key",
	}
	actual := tag.generate()
	expected := []byte("`json:\"key\"`")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestSample1(t *testing.T) {
	file, err := os.Open("../test_data/a.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gen := NewGenerator("test", "")
	if err := gen.ReadSchema(file, utils.FileName(file)); err != nil {
		panic(err)
	}
	actual, _ := gen.Generate(true)
	if len(actual) < 1 {
		t.Errorf("got %v\n", string(actual))
	}
}

func TestSample2(t *testing.T) {
	file, err := os.Open("../test_data/disk.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gen := NewGenerator("test", "")
	if err := gen.ReadSchema(file, utils.FileName(file)); err != nil {
		panic(err)
	}
	actual, _ := gen.Generate(false)
	if len(actual) < 1 {
		t.Errorf("got %v\n", string(actual))
	}
}

func TestSample3(t *testing.T) {
	file, err := os.Open("../test_data/qiita.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	gen := NewGenerator("test", "")
	if err := gen.ReadSchema(file, utils.FileName(file)); err != nil {
		panic(err)
	}
	actual, _ := gen.Generate(false)
	if len(actual) < 1 {
		t.Errorf("got %v\n", string(actual))
	}
}
