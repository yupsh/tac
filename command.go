package command

import (
	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Tac(parameters ...any) yup.Command {
	return command(yup.Initialize[yup.File, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return yup.AccumulateAndProcess(func(lines []string) []string {
		// Reverse all lines
		reversed := make([]string, len(lines))
		for i, line := range lines {
			reversed[len(lines)-1-i] = line
		}
		return reversed
	}).Executor()
}
