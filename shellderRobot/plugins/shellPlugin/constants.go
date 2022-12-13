package shellPlugin

const (
	ShellToUseUnix = "bash"
	ShellToUseWin  = "cmd"
)

const (
	vServersCmd             = "vservers"
	downloadCmd             = "download"
	uploadCmd               = "upload"
	dlCmd                   = "dl"
	ulCmd                   = "ul"
	exitCmd                 = "exit"
	executeCancelDataPrefix = "caEx"
	cbDataSep               = "_"
)

const (
	CommandExecuteTypeCmd = 0 + iota
	CommandExecuteTypePowerShell
)

// const (
// unsupportedMessage = "Unsupported operating system: " + runtime.GOOS
// )
