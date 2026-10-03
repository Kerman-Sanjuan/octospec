//go:build !windows

package session

import (
	"os"
	"syscall"
)

// processAlive reports whether pid is a live process, using signal 0, which
// probes liveness without sending a signal.
func processAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}
