// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// helpmux.go [created: Sun, 30 Jun 2013]

/*
Package helpmux provides 'help' functionality for complex command line programs.
*/
package helpmux

import (
	"fmt"
	"os"

	"github.com/bmatsuo/cmdmux"
)

var defaultHelpMux = New()

func Register(name string, cmd cmdmux.Command) error {
	return defaultHelpMux.Register(name, cmd)
}

func RegisterFunc(name string, cmd cmdmux.CommandFunc) error {
	return defaultHelpMux.RegisterFunc(name, cmd)
}

func RegisterHelp(name string, cmd cmdmux.Command) error {
	return defaultHelpMux.RegisterHelp(name, cmd)
}

func RegisterHelpFunc(name string, cmd cmdmux.CommandFunc) error {
	return defaultHelpMux.RegisterHelpFunc(name, cmd)
}

func Exec(name string, args []string) {
	defaultHelpMux.Exec(name, args)
}

func Help(name string, args []string) {
	defaultHelpMux.Help(name, args)
}

type HelpCommand interface {
	cmdmux.Command
	Help(name string, args []string)
}

type HelpMux struct {
	CmdMissing  func(name string, args []string)
	CmdUnknown  func(name string, args []string)
	HelpDefault func(name string, args []string)
	cmd         *cmdmux.Mux
	help        *cmdmux.Mux
}

func New() *HelpMux {
	mux := new(HelpMux)

	mux.help = cmdmux.NewMux()
	mux.help.CmdUnknown = func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: unknown help topic\n", args[0])
	}
	mux.help.CmdMissing = func(name string, args []string) {
		mux.HelpDefault(name, args)
	}

	mux.cmd = cmdmux.NewMux()
	mux.cmd.CmdUnknown = func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: unknown command\n", args[0])
	}
	mux.cmd.CmdMissing = func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: missing command\n", name)
	}

	return mux
}

// execute a help topic specified in the first argument.
func (mux *HelpMux) Help(name string, args []string) {
	if len(args) == 0 {
		if mux.HelpDefault != nil {
			mux.HelpDefault(name, args)
			return
		}
		fmt.Fprintf(os.Stderr, "%s: no help\n", name)
	} else {
		mux.help.Exec(name+" "+args[0], args[1:])
	}
}

// execute a command specified in the first argument.
func (mux *HelpMux) Exec(name string, args []string) {
	if len(args) == 0 {
		mux.CmdMissing(name, args)
		return
	}
	name = fmt.Sprintf("%s %s", name, args[0])
	args = args[1:]
	if args[0] == "help" {
		if len(args) == 1 {
			mux.Help(name, args[1:])
			return
		}
	}
	mux.cmd.Exec(name, args)
}

// register a Command with the given name. if the command is a HelpCommand,
// then a help topic is registered with the same name
func (mux *HelpMux) Register(name string, cmd cmdmux.Command) error {
	if cmd, ok := cmd.(HelpCommand); ok {
		err := mux.RegisterHelpFunc(name, cmd.Help)
		if err != nil {
			return err
		}
	}
	return mux.cmd.Register(name, cmd)
}

// register a command function. no help topic will be registered.
func (mux *HelpMux) RegisterFunc(name string, cmd cmdmux.CommandFunc) error {
	return mux.Register(name, cmd)
}

// register a help topic.
func (mux *HelpMux) RegisterHelp(name string, cmd cmdmux.Command) error {
	return mux.help.Register(name, cmd)
}

// register a help topic function.
func (mux *HelpMux) RegisterHelpFunc(name string, cmd cmdmux.CommandFunc) error {
	return mux.help.Register(name, cmd)
}
