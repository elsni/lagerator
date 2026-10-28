package testhome

import (
	"os"
	"testing"
)

// TestMain verifies test execution uses the same isolated home on every platform.
func TestMain(m *testing.M) { os.Exit(Run(m)) }

// TestIsolation checks both platform variables and the operating system's selected home.
func TestIsolation(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if home == "" || home != os.Getenv("HOME") || home != os.Getenv("USERPROFILE") {
		t.Fatalf("home is not isolated: %q", home)
	}
	env := Environment("isolated-child")
	if env[len(env)-3] != "HOME=isolated-child" || env[len(env)-2] != "USERPROFILE=isolated-child" {
		t.Fatal("child environment is not isolated")
	}
}
