package utils

func LDFlagsCheck(args []string, versionFn func(), helpFn func()) {
	if len(args) == 2 &&
		(args[1] == "--version" ||
			args[1] == "-version" ||
			args[1] == "-v" ||
			args[1] == "version") {
		versionFn()
	}

	if len(args) == 2 &&
		(args[1] == "--help" ||
			args[1] == "-help" ||
			args[1] == "-h" ||
			args[1] == "help") {
		helpFn()
	}
}
