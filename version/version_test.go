// Copyright 2017 aaharu All rights reserved.
// This source code is licensed under the BSD-style license found in
// the LICENSE file in the root directory of this source tree.

package version

import "testing"

func TestString(t *testing.T) {
	actual := String()
	expected := "schemarshal 0.6.0"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
