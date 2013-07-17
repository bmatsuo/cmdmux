// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cmdmux.go [created: Sun, 30 Jun 2013]

/*
subcommands
*/
package cmdmux

var defaultMux = NewMux()

// call Exec() on the default mux.
func Exec(name string, args []string) {
	defaultMux.Exec(name, args)
}

// register a function with the default mux.
func RegisterFunc(name string, cmd CommandFunc) error {
	return defaultMux.RegisterFunc(name, cmd)
}

// register a Command with the default mux.
func Register(name string, cmd Command) error {
	return defaultMux.Register(name, cmd)
}

// retrieve the names of commands registered with the default mux.
func Commands() []string {
	return defaultMux.Commands()
}

// a sub command. the Exec() method may or may not call os.Exit().
type Command interface {
	Exec(name string, args []string)
}

// an function that implements Command.
type CommandFunc func(string, []string)

func (fn CommandFunc) Exec(name string, args []string) {
	fn(name, args)
}
