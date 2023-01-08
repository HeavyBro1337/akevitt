package input

import "strings"

const COMMAND_PREFIX string = "/"

type InputType int

const (
	Ignore InputType = iota
	Command
	Message
)

// Entry point for all client input
func ParseInput(inp string) (status InputType) {

	// Check that the string is empty, otherwise see if its q/Q
	if len(inp) == 0 {
		return Ignore
	} else if strings.HasPrefix(inp, COMMAND_PREFIX) {
		// We entered command
		return Command
	}
	return Message
}
