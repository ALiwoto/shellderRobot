package shellPlugin

import "runtime"

const (
	ShellToUseUnix = "bash"
	ShellToUseWin  = "cmd"
)

const (
	downloadCmd = "download"
	uploadCmd   = "dl"
	dlCmd       = "download"
	ulCmd       = "ul"
	exitCmd     = "exit"
)

const (
	unsupportedMessage = "Unsupported operating system: " + runtime.GOOS
)
