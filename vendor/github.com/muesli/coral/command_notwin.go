//go:build !windows
// +build !windows

package coral

var preExecHookFn func(*Command)
