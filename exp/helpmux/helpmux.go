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
	"strings"

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

func Commands() []string {
	return defaultHelpMux.Commands()
}

func HelpTopics() []string {
	return defaultHelpMux.HelpTopics()
}

type HelpCommand interface {
	cmdmux.Command
	Help(name string, args []string)
}

type HelpMux struct {
	CmdMissing  cmdmux.Command
	CmdUnknown  cmdmux.Command
	HelpDefault cmdmux.Command
	cmd         *cmdmux.Mux
	help        *cmdmux.Mux
}

func New() *HelpMux {
	mux := new(HelpMux)

	mux.help = cmdmux.NewMux()
	mux.help.CmdUnknown = cmdmux.CommandFunc(func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: unknown help topic\n", args[0])
	})
	mux.help.CmdMissing = cmdmux.CommandFunc(func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: help topic missing\n", name, args)
	})

	mux.cmd = cmdmux.NewMux()
	mux.cmd.CmdUnknown = cmdmux.CommandFunc(func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: unknown command\n", args[0])
	})
	mux.cmd.CmdMissing = cmdmux.CommandFunc(func(name string, args []string) {
		fmt.Fprintf(os.Stderr, "%s: missing command\n", name)
	})

	return mux
}

// execute a help topic specified in the first argument.
func (mux *HelpMux) Help(name string, args []string) {
	fmt.Println(name, args)
	if len(args) == 0 {
		if mux.HelpDefault != nil {
			mux.HelpDefault.Exec(name, args)
			return
		} else {
			fmt.Fprintf(os.Stderr, "%s: help topics\n\t", name)
			fmt.Fprintf(os.Stderr, strings.Join(mux.HelpTopics(), "\n\t"))
			fmt.Fprintln(os.Stderr)
			os.Exit(1)
		}
	} else {
		mux.help.Exec(name, args)
	}
}

// execute a command specified in the first argument.
func (mux *HelpMux) Exec(name string, args []string) {
	if len(args) == 0 {
		mux.CmdMissing.Exec(name, args)
		return
	}
	base := args[0]
	name = fmt.Sprintf("%s %s", name, base)
	if base == "help" {
		args = args[1:]
		mux.Help(name, args)
		return
	}
	mux.cmd.Exec(name, args)
}

// registered command names.
func (mux *HelpMux) Commands() []string {
	return mux.cmd.Commands()
}

// registered help topics.
func (mux *HelpMux) HelpTopics() []string {
	return mux.help.Commands()
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
