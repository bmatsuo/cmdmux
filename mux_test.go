// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mux_test.go [created: Tue, 16 Jul 2013]

package cmdmux

import (
	"testing"
)

func TestMuxRegister(t *testing.T) {
	// i'm going to cheat for now and only test RegisterFunc
	mux := NewMux()
	noop := func(n string, a []string) {}

	for i, name := range []string{
		"abc", "def", "ghi",
		"adg", "beh", "cfi",
	} {
		err := mux.RegisterFunc(name, noop)
		if err != nil {
			t.Errorf("[%d] couldn't register %q; %v", i, name, err)
		}
	}

	dup := "def"
	err := mux.RegisterFunc(dup, noop) // double register
	if err == nil {
		t.Errorf("registered twice; %q", dup)
	}

	err = mux.RegisterFunc("xyz", nil)
	if err == nil {
		t.Errorf("registered nil func")
	}
	err = mux.Register("xyz", nil)
	if err == nil {
		t.Errorf("registered nil command")
	}

}
