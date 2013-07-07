// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// helpmux.go [created: Sun, 30 Jun 2013]

// Package helpmux does ....
package helpmux

import (
	"fmt"
	"os"

	"github.com/bmatsuo/cmdmux"
)

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
		mux.HelpDefault(name, args)
	}
	mux.help.CmdMissing = func(name string, args []string) {
		mux.HelpDefault(name, args)
	}

	mux.cmd = cmdmux.NewMux()
	mux.cmd.CmdUnknown = func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: unknown command\n", args[0])
		mux.HelpDefault(name, args)
	}
	mux.cmd.CmdMissing = func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: missing command\n", name)
		mux.HelpDefault(name, args)
	}

	return mux
}

func (mux *HelpMux) Exec(name string, args []string) {
	if len(args) == 0 {
		mux.CmdMissing(name, args)
		return
	}
	name = fmt.Sprintf("%s %s", name, args[0])
	args = args[1:]
	if args[0] == "help" {

		if len(args) == 1 {
			mux.HelpDefault(name, args)
			return
		}
	}
	mux.cmd.Exec(name, args)
}

func (mux *HelpMux) Register(name string, cmd cmdmux.Command) error {
	var help cmdmux.Command = cmdmux.CommandFunc(func(name string, args []string) {
		fmt.Printf("no help for `%s`", name)
	})
	if cmd, ok := cmd.(HelpCommand); ok {
		help = cmdmux.CommandFunc(cmd.Help)
	}
	err := mux.help.Register(name, help)
	if err != nil {
		return err
	}
	return mux.cmd.Register(name, help)
}

func (mux *HelpMux) RegisterFunc(name string, cmd cmdmux.CommandFunc) error {
	return mux.Register(name, cmd)
}

func (mux *HelpMux) RegisterHelp(name string, cmd cmdmux.Command) error {
	return mux.help.Register(name, cmd)
}

func (mux *HelpMux) RegisterHelpFunc(name string, cmd cmdmux.CommandFunc) error {
	return mux.help.Register(name, cmd)
}
