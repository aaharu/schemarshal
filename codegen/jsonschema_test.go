// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package codegen

import (
	"os"
	"testing"

	"github.com/aaharu/schemarshal/utils"
)

func TestSample1(t *testing.T) {
	file, err := os.Open("../test_data/a.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	js, err := ReadSchema(file)
	if err != nil {
		panic(err)
	}
	gen := NewGenerator("test", "")
	jsType, _ := js.parse(utils.UpperCamelCase(utils.FileName(file)), gen)
	actual := jsType.generate()
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
	js, err := ReadSchema(file)
	if err != nil {
		panic(err)
	}
	gen := NewGenerator("test", "")
	jsType, _ := js.parse(utils.UpperCamelCase(utils.FileName(file)), gen)
	actual := jsType.generate()
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
	js, err := ReadSchema(file)
	if err != nil {
		panic(err)
	}
	gen := NewGenerator("test", "")
	jsType, err := js.parse(utils.UpperCamelCase(utils.FileName(file)), gen)
	if err != nil {
		panic(err)
	}
	actual := jsType.generate()
	if len(actual) < 1 {
		t.Errorf("got %v\n", string(actual))
	}
}
