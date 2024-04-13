package process

/*

Tools to manage OS level processes

TODO
 * Convert to a class with member methods instead of a collection of loose functions
 * Make paths to external/shell tools configurable with sane defaults (bash, nohup, awk, echo, etc)
 * Convert to use Go built-ins when/where possible, aim for platform agnostic (currently *nix bound)
 * Add process watchdog that tracks PID for each process that we spawn
 * Add functionality to terminate a specified process (may as well use PID as the key; OS offers no protections here)

*/

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

// Run the specified command as a background process (*nix) and return its process ID (PID)
// NOTE: Platform-specific, uses bash, nohup, echo, awk  and /dev/null redirection on *nix hosts
func RunBackgroundCommand_Nix(command string) (int, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	bgCmd := exec.Command(fmt.Sprintf(
		"/bin/bash -c nohup %s < /dev/null &>/dev/null & echo -n $! | awk '/[0-9]+$/{ print $0 }'",
		command,
	))
	bgCmd.Stdout = &stdout
	bgCmd.Stderr = &stderr

	err := bgCmd.Run()
	if nil != err {
		return 0, errors.New(fmt.Sprintf("%v: %v", err.Error(), stderr.String()))
	}

	pid := strconv.Atoi(stdout.String())
	if 0 == pid { return 0, errors.New("Invalid PID! (shouldn't happen without error)") }

	return pid, nil
}

func RunCommand_Nix(command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// ref: https://stackoverflow.com/questions/1877045/how-do-you-get-the-output-of-a-system-command-in-go
	output, err := exec.Command(fmt.Sprintf( "/bin/bash %s", command)).Output()
	if nil != err {
		return "", errors.New(fmt.Sprintf("%v: %v", err.Error(), stderr.String()))
	}

	return string(pid), nil
}

