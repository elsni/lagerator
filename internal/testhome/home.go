// Package testhome isolates inventory files during tests on Unix and Windows.
package testhome

import (
	"os"
	"testing"
)

// Run executes a package's tests with a temporary platform-independent home directory.
func Run(m *testing.M) int {
	home, err := os.MkdirTemp("", "lgrttest-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(home)
	for _, key := range []string{"HOME", "USERPROFILE"} {
		if err := os.Setenv(key, home); err != nil {
			panic(err)
		}
	}
	return m.Run()
}

// Environment returns subprocess settings that isolate both Unix and Windows inventories.
func Environment(home string) []string {
	return append(os.Environ(), "HOME="+home, "USERPROFILE="+home, "TZ=UTC")
}
