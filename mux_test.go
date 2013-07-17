// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mux_test.go [created: Tue, 16 Jul 2013]

package cmdmux

import (
	"testing"
)

func noopFunc(name string, args []string) {}

var noop Command = CommandFunc(noopFunc)

func TestMuxRegister(t *testing.T) {
	mux := NewMux()
	funcmux := NewMux()

	// register success
	for i, name := range []string{
		"abc", "def", "ghi",
		"adg", "beh", "cfi",
	} {
		err := mux.Register(name, noop)
		if err != nil {
			t.Errorf("[%d] couldn't register %q; %v", i, name, err)
		}
		err = funcmux.RegisterFunc(name, noopFunc)
		if err != nil {
			t.Errorf("[%d] couldn't register func %q; %v", i, name, err)
		}
	}

	// double register
	dup := "def"
	err := mux.Register(dup, noop)
	if err == nil {
		t.Errorf("registered twice; %q", dup)
	}
	err = funcmux.RegisterFunc(dup, noopFunc)
	if err == nil {
		t.Errorf("registered func twice; %q", dup)
	}

	// nil register
	err = mux.RegisterFunc("xyz", nil)
	if err == nil {
		t.Errorf("registered nil func")
	}
	err = mux.Register("xyz", nil)
	if err == nil {
		t.Errorf("registered nil command")
	}

}
