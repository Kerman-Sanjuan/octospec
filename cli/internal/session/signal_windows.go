//go:build windows

package session

// processAlive is a stub on Windows; octospec does not ship a Windows binary.
func processAlive(pid int) bool { return false }
