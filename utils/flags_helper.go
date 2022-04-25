package utils

import "os"

func LDFlagsCheck(args []string, versionFn func(), helpFn func()) {
	if len(args) == 2 &&
		(args[1] == "--version" ||
			args[1] == "-version" ||
			args[1] == "-v" ||
			args[1] == "version") {
		versionFn()
		os.Exit(0)
	}

	if len(args) == 2 &&
		(args[1] == "--help" ||
			args[1] == "-help" ||
			args[1] == "-h" ||
			args[1] == "help") {
		helpFn()
		os.Exit(0)
	}
}
