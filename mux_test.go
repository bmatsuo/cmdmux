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

func TestMuxExec(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("panic;", e)
		}
	}()

	x := 0
	mux := NewMux()
	err := mux.RegisterFunc("blah", func(name string, args []string) { x++ })
	if err != nil {
		t.Fatalf("couldn't register test function;", err)
	}
	mux.Exec("bleh", []string{"blah"})
	if x == 0 {
		t.Fatal("command did not execute")
	}
}

func TestMuxExecUnknown(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("panic;", e)
		}
	}()
	mux := NewMux()
	caught := false
	mux.CmdUnknown = CommandFunc(func(string, []string) {
		caught = true
	})
	mux.Exec("bleh", []string{"meh"})
	if !caught {
		t.Fatal("error not caught")
	}
}

func TestMuxExecMissing(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("panic;", e)
		}
	}()
	mux := NewMux()
	caught := false
	mux.CmdMissing = CommandFunc(func(string, []string) {
		caught = true
	})
	mux.Exec("bleh", []string{})
	if !caught {
		t.Fatal("error not caught")
	}
}

func TestMuxCommands(t *testing.T) {
	mux := NewMux()
	commands := []string{
		"123", "abc", "456",
		"def", "789", "ghi",
	}
	for _, name := range commands {
		mux.Register(name, noop)
	}
	_commands := mux.Commands()
	if len(commands) != len(_commands) {
		t.Errorf("commands don't match %d != %d", len(commands), len(_commands))
	} else {
		for i := range commands {
			if commands[i] != _commands[i] {
				t.Errorf("command %d doesn't match %q != %q", i, commands[i], _commands[i])
			}
		}
	}
}
