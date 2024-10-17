// Package command provides utilities to manage and execute shell commands.
// It supports executing individual commands, combining multiple commands into one, and running
// those commands within a specified environment. Additionally, the package offers methods for
// initializing, executing, and handling installation of commands, often in conjunction with
// external tools like mvdan/sh for shell script parsing and execution.
//
// The Commands type represents a collection of shell commands that can be processed together.
// The Command type encapsulates a single shell command, with methods to manipulate, format, and
// execute it.
package command
