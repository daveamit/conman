package conman

import "os"

// ExitCode represents all possible error codes logged for ref.
type ExitCode int

const (
	// SuccessExitCode indicates success
	SuccessExitCode ExitCode = iota
	// FailedToCreateLogsExitCode indicates was unable to instantiate logs.
	FailedToCreateLogsExitCode
)

func exit(code ExitCode) {
	os.Exit(int(code))
}
