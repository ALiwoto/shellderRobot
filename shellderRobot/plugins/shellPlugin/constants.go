package shellPlugin

import "runtime"

const (
	ShellToUseUnix = "bash"
	ShellToUseWin  = "cmd"
)

const (
	downloadCmd = "download"
	uploadCmd   = "upload"
	dlCmd       = "dl"
	ulCmd       = "ul"
	exitCmd     = "exit"
)

const (
	unsupportedMessage = "Unsupported operating system: " + runtime.GOOS
)
