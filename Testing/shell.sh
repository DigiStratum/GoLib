#!/bin/bash

# Reset debug logging
debugLog=debug.log
if [ -f "$debugLog" ]; then
        rm -f $debugLog
fi

#-------------------------------------------------------------------------------
# General reusable functions
#-------------------------------------------------------------------------------
# ref: https://stackoverflow.com/questions/9893667/is-there-a-way-to-write-a-bash-function-which-aborts-the-whole-execution-no-mat
trap "exit 1" TERM
export TOP_PID=$$

# Use this to echo to the console AND log the output
say() {
	local msg="$1"
	echo -e "$msg" >> $debugLog
	echo -e "$msg" >&2
}

# Use this to die with a message
die() {
	say "$1"
	# We need this to exit ths script, not just the function
	kill -s TERM $TOP_PID
}

# Use this to check the result of system/exec commands that are expected to return 0 for SUCCESS (and die)
die_on_error() {
	local actual="$1"
	local msg="$2"
	if [ ${#actual} -ne 0 ] && [ "$actual" != "0" ]; then
		die "$msg"
	fi
}

# Use this to check the result of function call that are expected to return 0 for ERROR (and die)
check() {
	local actual="$1"
	local msg="$2"
	if [ ${#actual} -eq 0 ]; then
		die "$msg"
	fi
}

# Use this to check any explicit expected value against actual (and die)
expect() {
	local actual="$1"
	local expect="$2"
	local msg="$3"
	if [ "$actual" != "$expect" ]; then
		die "$msg"
	fi
}

