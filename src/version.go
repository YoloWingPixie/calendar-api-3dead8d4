package main

import (
	"fmt"
	"runtime"
)

// Version information - these will be set at build time via -ldflags
var (
	Version   = "dev"     // Project version (e.g., v1.0.0)
	Commit    = "unknown" // Git commit hash
	Date      = "unknown" // Build date
	GoVersion = runtime.Version()
)

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	return fmt.Sprintf("Calendar API %s (commit: %s, built: %s, go: %s)",
		Version, Commit, Date, GoVersion)
}

// GetVersion returns just the version string
func GetVersion() string {
	return Version
}
