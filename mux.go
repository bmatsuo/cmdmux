// Copyright 2013, Bryan Matsuo. All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mux.go [created: Sun, 30 Jun 2013]

package cmdmux

import (
	"errors"
	"fmt"
)

var ErrNilRegister = errors.New("cannot register nil command")
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

// a subcommand mux. a simple name lookup table.
type Mux struct {
	CmdMissing Command
	CmdUnknown Command
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
	mux.CmdMissing.Exec(name, args)
}

func (mux *Mux) cmdUnknown(name string, args []string) {
	if mux.CmdUnknown == nil {
		defaultCmdUnknown(name, args)
		return
	}
	mux.CmdUnknown.Exec(name, args)
}

// implements the Command interaface. execute the command named in args[0].
// if there is no such command, mux.CmdUnknown() is called. if the argument
// list is empty, mux.CmdMissing() is called. a runtime panic if the command
// is missing or unrecognized and the appropriate command is nil.
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

// register a function. see Register.
func (mux *Mux) RegisterFunc(name string, cmd CommandFunc) error {
	if cmd == nil {
		return ErrNilRegister
	}
	return mux.Register(name, cmd)
}

// register a command with a given name. the name cannot already be taken and
// the command cannot be nil.
func (mux *Mux) Register(name string, cmd Command) error {
	if cmd == nil {
		return ErrNilRegister
	}
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

// returns the names of registered commands in the order they were registered.
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
