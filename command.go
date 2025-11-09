package command

import (
	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Tac(parameters ...any) gloo.Command {
	return command(gloo.Initialize[gloo.File, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return gloo.AccumulateAndProcess(func(lines []string) []string {
		// Reverse all lines
		reversed := make([]string, len(lines))
		for i, line := range lines {
			reversed[len(lines)-1-i] = line
		}
		return reversed
	}).Executor()
}
