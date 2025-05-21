#!/bin/bash


# TODO:
# * Stop using die/exit/TERM trap, etc. just return exit codes through the function calls back to parent, parent exits as needed.

# Which directory should we return to upon death? Defaults to no change on exit (it stays where it dies)
DEATH_DIR="$PWD"

# Reset debug logging
#debugLog="debug.log"
debugLog=""
if [ -f "$debugLog" ]; then
        rm -f $debugLog
fi

if [ "$debugLog" != "" ]; then
        touch $debugLog
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
  if [ -f "$debugLog" ]; then
	  echo -e "$msg" >> $debugLog
  fi
	echo -e "$msg" >&2
}

# Set the directory that death will return us to (should just be called once by the parent script)
set_death_dir() {
	DEATH_DIR="$1"
}

# Use this to die with a message
die() {
	say "$1"
	# We need this to exit the script, not just the function
	cd $DEATH_DIR
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


#-------------------------------------------------------------------------------
# Test Runners
#-------------------------------------------------------------------------------

# Run tests, no coverage report
runtests() {
	local TESTRUNNER_DIR="$1"
	local RUN_TEST=""

	if [ "$#" -eq 2 ]; then
		RUN_TEST="-run $2"
	fi

	say "BEGIN: $TESTRUNNER_DIR"
	cd $TESTRUNNER_DIR

	# Every go source file must have a package name matching that of the directory
	local EXPECTPKGNAME=`basename $TESTRUNNER_DIR`
	say "\nPACKAGE: $EXPECTPKGNAME"
	for f in `ls -1 $TESTRUNNER_DIR/*.go`; do
		local PKGMATCH=`grep -i "^package $EXPECTPKGNAME$" $f | wc -l`
		if [ "$PKGMATCH" -eq 0 ]; then
			say "\t\tWARN: Expected package name missing from source file ('$f')"
		fi
	done

	# Run tests
	#go test -v -count=1 -covermode=atomic -coverprofile=test_coverage.txt $RUN_TEST *.go
	output=`go test -v -count=1 -covermode=atomic -coverprofile=.test_coverage.txt $RUN_TEST *.go`
	res=$?
	echo "$output" | grep -v "^--- PASS:" | grep -v "^--- FAIL:"
	if [ 0 != $res ]; then
		say "^^^^ ERROR!"
		return 1
	fi

	say "\nDONE: $TESTRUNNER_DIR"
}

