#! /bin/bash

# Author: Parker Jones
# Description: This is the build/installer script for linux for the sheep-form program
#		The sheep-form program is used for creating form boilerplate in various languages from simple
#   quick to write, minimalistic templates
BIN_DIR=~/.local/bin
BASE_DIR=~/.local/share/sheep/sheep-form
GO_FILES=$(find ./ -maxdepth 1 -type f -name "src/*.go")

case $1 in
	build)
		go build -o ./sheep-forms $GO_FILES
		;;
	install)
		# Build the executable
		go build -o ./sheep-forms $GO_FILES
		# Create and move template dir
		mkdir --parent $BASE_DIR; cp -r ./templates $_
		# Move executable to .local/bin
		cp ./sheep-forms $BIN_DIR/sheep-forms
		;;
	run)
		./sheep-forms
		;;
	*)
		echo "Unrecognized argument. Valid args are 'build', 'run', and 'install'"
		;;
esac
