package version

// Version is the current version of the dev CLI tool.
// It can be overridden at build time using ldflags:
//
//	go build -ldflags "-X dev/internal/version.Version=v1.2.0" .
var Version = "v0.1.0"
