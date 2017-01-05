// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package schemarshal

import "testing"

func TestMarshalTagEmpty(t *testing.T) {
	actual := MarshalTag("address", false)
	expected := "`json:\"address,omitempty\"`"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestMarshalTag(t *testing.T) {
	actual := MarshalTag("address", true)
	expected := "`json:\"address\"`"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestUcfirst(t *testing.T) {
	actual := Ucfirst("address")
	expected := "Address"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
