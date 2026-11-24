package install_test

import (
	"archive/tar"
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/path/file"
)

// TestUpdateRunningExecutable verifies that a running Linux process keeps its image across an update.
//
//nolint:paralleltest // Do not fork while another subtest holds a writable executable descriptor.
func TestUpdateRunningExecutable(t *testing.T) {
	// Keep the shell's name: BusyBox selects its applet from argv[0].
	const exe = "sh"

	shell, err := file.New(exe).Which()
	if err != nil {
		t.Fatal(err)
	}

	old, err := shell.Read()
	if err != nil {
		t.Fatal(err)
	}

	const replacement = "#!/bin/sh\necho new\n"

	for _, tc := range []struct {
		path string
		mode mode.Mode
		perm os.FileMode
	}{
		{exe, mode.Extract, 0o700},        // Raw downloads retain the existing executable permissions.
		{"tool.tar", mode.Extract, 0o750}, // Archive members supply their own permissions.
		{"tool.tar", mode.Find, 0o750},
	} {
		t.Run(tc.mode.String()+"/"+tc.path, func(t *testing.T) {
			out := t.TempDir()
			target := filepath.Join(out, exe)

			if err := file.New(target).Write(old, 0o700); err != nil {
				t.Fatal(err)
			}

			const script = `echo ready; read -r line && echo "${line}"`

			cmd := exec.CommandContext(t.Context(), target, "-c", script) //nolint:gosec // Local test shell.

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatal(err)
			}

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}

			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				_ = stdin.Close()

				if cmd.ProcessState == nil {
					_ = cmd.Wait()
				}
			})

			reader := bufio.NewReader(stdout)
			if line, err := reader.ReadString('\n'); err != nil || line != "ready\n" {
				t.Fatalf("old process did not become ready: %q, %v", line, err)
			}

			body := []byte(replacement)

			if tc.path != exe {
				body = testTar(t, replacement, tar.Header{Name: exe, Mode: 0o750})
			}

			d := testData(t, out, tc.path, body)

			d.Mode, d.Exe, d.Patterns = tc.mode.String(), exe, []string{exe}

			if _, err := install.Download(d); err != nil {
				t.Fatal(err)
			}

			want := tc.perm

			if tc.mode == mode.Find {
				// Find mode creates the executable subject to the process umask.
				ref := file.New(t.TempDir(), "reference")
				if err := ref.Write(nil, want); err != nil {
					t.Fatal(err)
				}

				info, err := ref.Info()
				if err != nil {
					t.Fatal(err)
				}

				want = info.Mode().Perm()
			}

			expectFile(t, target, replacement, want)

			if _, err := io.WriteString(stdin, "still running\n"); err != nil {
				t.Fatal(err)
			}

			if err := stdin.Close(); err != nil {
				t.Fatal(err)
			}

			if line, err := reader.ReadString('\n'); err != nil || line != "still running\n" {
				t.Fatalf("old process: %q, %v", line, err)
			}

			if err := cmd.Wait(); err != nil {
				t.Fatal(err)
			}

			got, err := exec.CommandContext(t.Context(), target).Output() //nolint:gosec // Local test fixture.
			if err != nil || string(got) != "new\n" {
				t.Fatalf("new process: %q, %v", got, err)
			}
		})
	}
}
