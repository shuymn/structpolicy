package ptrstruct_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_HelpOmitsFixFlags(t *testing.T) {
	t.Parallel()

	binPath := buildBinary(t)
	help := exec.Command(binPath, "-help")
	out, err := help.CombinedOutput()
	if err == nil {
		t.Fatalf("expected help command to exit non-zero, got success\n%s", out)
	}

	text := string(out)
	for _, forbidden := range []string{"-fix", "-diff"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("help output should not mention %s\n%s", forbidden, text)
		}
	}
}

func TestCLI_RejectsFixFlag(t *testing.T) {
	t.Parallel()

	binPath := buildBinary(t)
	cmd := exec.Command(binPath, "-fix", "./...")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected -fix to be rejected\n%s", out)
	}
	if !strings.Contains(string(out), "flag provided but not defined: -fix") {
		t.Fatalf("unexpected output\n%s", out)
	}
}

func TestCLI_ReportsDiagnostics(t *testing.T) {
	t.Parallel()

	binPath := buildBinary(t)
	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "go.mod"), "module example.com/cli\n\ngo 1.25.0\n")
	writeFile(t, filepath.Join(workspace, "cli.go"), `package cli

type User struct {
	Name string
}

func (u User) Normalize() {}

func Save(users []User) {}
`)

	cmd := exec.Command(binPath, "./...")
	cmd.Dir = workspace
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected diagnostics exit\n%s", out)
	}

	text := string(out)
	if !strings.Contains(text, "receiver uses value struct User; use *User") {
		t.Fatalf("output missing receiver diagnostic\n%s", text)
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()

	repoRoot, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(t.TempDir(), "ptrstruct")
	build := exec.Command("go", "build", "-o", binPath, "./cmd/ptrstruct")
	build.Dir = repoRoot
	if buildOut, buildErr := build.CombinedOutput(); buildErr != nil {
		t.Fatalf("build ptrstruct: %v\n%s", buildErr, buildOut)
	}

	return binPath
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
