package input

import "strings"

const MESSAGE_PREFIX string = "say"

type InputType int

const (
	Ignore InputType = iota
	Command
	Message
)

// Entry point for all client input
func ParseInput(inp string) (status InputType, parsedInput string) {

	// Check that the string is empty, otherwise see if its q/Q
	if len(inp) == 0 {
		return Ignore, inp
	} else if strings.HasPrefix(inp, MESSAGE_PREFIX) {
		// We entered command
		return Message, strings.Replace(inp, MESSAGE_PREFIX, "", 1)
	}
	return Command, strings.Replace(inp, MESSAGE_PREFIX, "", 1)
}
