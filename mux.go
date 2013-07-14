// Copyright 2013, Bryan Matsuo. All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mux.go [created: Sun, 30 Jun 2013]

package cmdmux

import (
	"errors"
	"fmt"
)

var ErrDoubleRegister = errors.New("name already registered")
var errCmdMissing = errors.New("no command given")

type cmdUnknownError string

func (err cmdUnknownError) String() string {
	return fmt.Sprint("unknown command:", string(err))
}

var defaultCmdMissing = func(name string, args []string) {
	panic(errCmdMissing)
}

var defaultCmdUnknown = func(name string, args []string) {
	panic(cmdUnknownError(name))
}

// a command mux.
type Mux struct {
	CmdMissing func(name string, args []string)
	CmdUnknown func(name string, args []string)
	cmdnames   []string
	table      commandTable
}

// create a new mux
func NewMux() *Mux {
	return &Mux{table: newCommandTable(5)}
}

func (mux *Mux) cmdMissing(name string, args []string) {
	if mux.CmdMissing == nil {
		defaultCmdMissing(name, args)
		return
	}
	mux.CmdMissing(name, args)
}

func (mux *Mux) cmdUnknown(name string, args []string) {
	if mux.CmdUnknown == nil {
		defaultCmdUnknown(name, args)
		return
	}
	mux.CmdUnknown(name, args)
}

// implements the Command interaface. execute the registered command named in
// args[0].
func (mux *Mux) Exec(name string, args []string) {
	if len(args) == 0 {
		mux.cmdMissing(name, args)
		return
	}
	cmd := mux.table.Read(args[0])
	if cmd == nil {
		mux.cmdUnknown(name, args)
		return
	}
	name = fmt.Sprintf("%s %s", name, args[0])
	cmd.Exec(name, args[1:])
}

// register a function
func (mux *Mux) RegisterFunc(name string, cmd CommandFunc) error {
	return mux.Register(name, cmd)
}

// register a command
func (mux *Mux) Register(name string, cmd Command) error {
	var err error
	mux.table.Write(name, func(_cmd Command) (Command, error) {
		if _cmd != nil {
			err = ErrDoubleRegister
			return nil, err
		}
		// update cmdnames inside the lock
		mux.cmdnames = append(mux.cmdnames, name)
		return cmd, nil
	})
	return err
}

// returns the names of registered commands.
func (mux *Mux) Commands() []string {
	return mux.cmdnames
}

// a threadsafe backing to a Mux
type commandTable chan map[string]Command

var errNoop = errors.New("noop")

func newCommandTable(size uint) commandTable {
	ch := make(commandTable, 1)
	ch <- make(map[string]Command, size)
	return ch
}

func (table commandTable) Commands(name string) []string {
	names := make([]string, 0, 5)
	t := <-table
	defer func() { table <- t }() // not really necessary
	for name := range t {
		names = append(names, name)
	}
	return names
}

func (table commandTable) Read(name string) Command {
	t := <-table
	defer func() { table <- t }()
	return t[name]
}

func (table commandTable) Write(name string, fn func(Command) (Command, error)) error {
	t := <-table
	defer func() { table <- t }()

	cmd, err := fn(t[name])
	switch err {
	case nil:
		break
	case errNoop:
		return nil
	default:
		return err
	}

	if cmd == nil {
		delete(t, name)
	} else {
		t[name] = cmd
	}
	return nil
}
