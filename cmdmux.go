// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cmdmux.go [created: Sun, 30 Jun 2013]

/*
command line sub command nano framework
*/
package cmdmux

var defaultMux = NewMux()

func Exec(name string, args []string) {
	defaultMux.Exec(name, args)
}

func RegisterFunc(name string, cmd CommandFunc) error {
	return defaultMux.RegisterFunc(name, cmd)
}

func Register(name string, cmd Command) error {
	return defaultMux.Register(name, cmd)
}

func Commands() []string {
	return defaultMux.Commands()
}

type Command interface {
	Exec(name string, args []string)
}

type CommandFunc func(string, []string)

func (fn CommandFunc) Exec(name string, args []string) {
	fn(name, args)
}
