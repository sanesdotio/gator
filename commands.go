package main

import "fmt"

type command struct {
	Name string
	args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (cmds *commands) run(state *state, cmd command) error {
	if command, exists := cmds.Commands[cmd.Name]; exists {
		return command(state, cmd)
	} else {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
}

func (cmds *commands) register(name string, function func(*state, command) error) error {
	if _, exists := cmds.Commands[name]; exists {
		return fmt.Errorf("command %s already exists", name)
	}
	cmds.Commands[name] = function
	return nil
}