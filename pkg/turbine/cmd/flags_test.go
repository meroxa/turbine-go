package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_execPath(t *testing.T) {
	path := execPath()
	if path == "" {
		t.Fatalf("path cannot be empty")
	}
}

type testString string

func (t testString) Set(s string) error { return nil }
func (t testString) String() string     { return string(t) }

func Test_requiredFlags(t *testing.T) {
	var (
		empty      = &flag.Flag{Name: "flag", Value: testString("")}
		valid      = &flag.Flag{Name: "flag", Value: testString("valid")}
		f          = &fakeExiter{}
		origExiter = exiter
	)
	exiter = f.Exit

	tests := []struct {
		name string
		fn   func(f *flag.Flag)
	}{
		{"required server flag", requiredServerFlag},
		{"required build flag", requiredBuildFlag},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(empty)
			if want := 2; f.rc != want {
				t.Fatalf("empty flag: got %d, want %d", f.rc, want)
			}
			f.rc = 0

			tc.fn(valid)
			if want := 0; f.rc != want {
				t.Fatalf("valid flag: got %d, want %d", f.rc, want)
			}
			f.rc = 0
		})
	}
	exiter = origExiter
}

func Test_parseFlags(t *testing.T) {
	osArgs := os.Args

	tests := []struct {
		name     string
		args     []string
		matcher  func() string
		wantCode int
	}{
		{
			name: "parses server flags",
			args: []string{
				"cmd", "build",
				"--gitsha", "1234",
				"--turbine-core-server", ":5151",
			},
			matcher: func() string {
				msg := []string{}
				if want := "1234"; buildGitSHA != want {
					msg = append(msg, fmt.Sprintf("buildGitSHA: got %s, want %s", buildGitSHA, want))
				}
				if want := ":5151"; buildListenAddr != want {
					msg = append(msg, fmt.Sprintf("buildListenAddr: got %s, want %s", buildListenAddr, want))
				}
				return strings.Join(msg, ", ")
			},
		},
		{
			name: "parses build flags",
			args: []string{
				"command", "server",
				"--serve-addr", ":5151",
				"--serve", "my-func",
			},
			matcher: func() string {
				msg := []string{}
				if want := "my-func"; serverFuncName != want {
					msg = append(msg, fmt.Sprintf("serverFuncName: got %s, want %s", serverFuncName, want))
				}
				if want := ":5151"; serverListenAddr != want {
					msg = append(msg, fmt.Sprintf("serverListenAddr: got %s, want %s", serverListenAddr, want))
				}
				return strings.Join(msg, ", ")
			},
		},
		{
			name: "usage and exit without subcommand",
			args: []string{
				"command", "foobar",
			},
			wantCode: 2,
			matcher:  func() string { return "" },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeExiter{}
			origExiter := exiter
			exiter = f.Exit
			os.Args = tc.args

			parseFlags()
			if got := tc.matcher(); got != "" {
				t.Fatalf("failed to match flags: %s", got)
			}
			if got := f.rc; got != tc.wantCode {
				t.Fatalf("got %d, want %d", got, tc.wantCode)
			}

			exiter = origExiter
		})
	}
	os.Args = osArgs
}

type fakeExiter struct {
	rc int
}

func (f *fakeExiter) Exit(code int) {
	f.rc = code
}
