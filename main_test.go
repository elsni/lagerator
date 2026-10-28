package main

import (
	"bytes"
	"encoding/json"
	"github.com/elsni/lagerator/internal/testhome"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestCLIOutputUnchanged compares terminal output with snapshots from before extraction.
func TestCLIOutputUnchanged(t *testing.T) {
	var cases []struct {
		Args   []string `json:"args"`
		Stdout string   `json:"stdout"`
		Stderr string   `json:"stderr"`
	}
	golden, err := os.ReadFile("testdata/cli-output.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(golden, &cases); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile("testdata/inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "lgrt")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.Args, " "), func(t *testing.T) {
			home := t.TempDir()
			if err := os.Mkdir(filepath.Join(home, ".lgrt"), 0700); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(home, ".lgrt", "lgrtdata.json")
			if err := os.WriteFile(path, fixture, 0600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(binary, tc.Args...)
			cmd.Env = testhome.Environment(home)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("run: %v\n%s", err, stderr.String())
			}
			// Windows uses its system timezone rather than the subprocess TZ setting.
			if runtime.GOOS == "windows" {
				for _, timestamp := range []int64{1700000000, 1700000060} {
					value := time.Unix(timestamp, 0)
					tc.Stdout = strings.ReplaceAll(tc.Stdout, value.UTC().Format("02.01.2006 15:04"), value.Local().Format("02.01.2006 15:04"))
				}
			}
			if stdout.String() != tc.Stdout || stderr.String() != tc.Stderr {
				t.Errorf("output changed\nstdout: %q\nwant:   %q\nstderr: %q\nwant:   %q", stdout.String(), tc.Stdout, stderr.String(), tc.Stderr)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, fixture) {
				t.Error("read-only command changed inventory data")
			}
		})
	}
}
