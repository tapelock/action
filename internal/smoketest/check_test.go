package smoketest

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// scriptPath locates scripts/check.sh relative to this test file, so the
// test works regardless of the working directory `go test` is run from.
func scriptPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	return filepath.Join(repoRoot, "scripts", "check.sh")
}

// writeFakeTapelock creates a fake `tapelock` executable in dir, found via
// PATH by check.sh: `tapelock check --cassette X ...` exits 1 if X's name
// contains "fail", 0 otherwise, and prints the cassette path it checked.
func writeFakeTapelock(t *testing.T, dir string) {
	t.Helper()
	script := `#!/usr/bin/env bash
if [ "$1" = "check" ]; then
  shift
  cassette=""
  while [ $# -gt 0 ]; do
    case "$1" in
      --cassette) cassette="$2"; shift 2 ;;
      *) shift ;;
    esac
  done
  echo "fake-tapelock check $cassette"
  case "$cassette" in
    *fail*) exit 1 ;;
    *) exit 0 ;;
  esac
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(dir, "tapelock"), []byte(script), 0o755); err != nil {
		t.Fatalf("write fake tapelock: %v", err)
	}
}

func runCheckScript(t *testing.T, binDir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(scriptPath(t), args...)
	cmd.Env = append(os.Environ(), "PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func writeFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestCheckScriptSingleFilePass(t *testing.T) {
	binDir := t.TempDir()
	writeFakeTapelock(t, binDir)

	cassette := filepath.Join(t.TempDir(), "ok.jsonl")
	writeFile(t, cassette, "{}\n", 0o644)

	if out, err := runCheckScript(t, binDir, cassette); err != nil {
		t.Fatalf("check.sh failed: %v\n%s", err, out)
	}
}

func TestCheckScriptSingleFileFail(t *testing.T) {
	binDir := t.TempDir()
	writeFakeTapelock(t, binDir)

	cassette := filepath.Join(t.TempDir(), "fail.jsonl")
	writeFile(t, cassette, "{}\n", 0o644)

	if _, err := runCheckScript(t, binDir, cassette); err == nil {
		t.Fatal("check.sh: want a non-zero exit for a failing cassette, got nil error")
	}
}

func TestCheckScriptDirectoryFailsIfAnyFileFails(t *testing.T) {
	binDir := t.TempDir()
	writeFakeTapelock(t, binDir)

	cassetteDir := t.TempDir()
	writeFile(t, filepath.Join(cassetteDir, "ok.jsonl"), "{}\n", 0o644)
	writeFile(t, filepath.Join(cassetteDir, "fail.jsonl"), "{}\n", 0o644)
	writeFile(t, filepath.Join(cassetteDir, "not-a-cassette.txt"), "ignored", 0o644)

	out, err := runCheckScript(t, binDir, cassetteDir)
	if err == nil {
		t.Fatalf("check.sh: want a non-zero exit when one file in the directory fails, got nil\n%s", out)
	}
	if !strings.Contains(out, "ok.jsonl") || !strings.Contains(out, "fail.jsonl") {
		t.Fatalf("check.sh did not check both cassette files:\n%s", out)
	}
	if strings.Contains(out, "not-a-cassette.txt") {
		t.Fatalf("check.sh checked a non-.jsonl file:\n%s", out)
	}
}

func TestCheckScriptDirectoryAllPass(t *testing.T) {
	binDir := t.TempDir()
	writeFakeTapelock(t, binDir)

	cassetteDir := t.TempDir()
	writeFile(t, filepath.Join(cassetteDir, "a.jsonl"), "{}\n", 0o644)
	writeFile(t, filepath.Join(cassetteDir, "b.jsonl"), "{}\n", 0o644)

	if out, err := runCheckScript(t, binDir, cassetteDir); err != nil {
		t.Fatalf("check.sh failed: %v\n%s", err, out)
	}
}

func TestCheckScriptEmptyDirectoryFails(t *testing.T) {
	binDir := t.TempDir()
	writeFakeTapelock(t, binDir)

	cassetteDir := t.TempDir() // no .jsonl files in it

	if _, err := runCheckScript(t, binDir, cassetteDir); err == nil {
		t.Fatal("check.sh: want a non-zero exit for a directory with no cassettes, got nil")
	}
}

func TestCheckScriptPassesThroughExtraArgs(t *testing.T) {
	binDir := t.TempDir()
	// A fake tapelock that just echoes every arg it received.
	writeFile(t, filepath.Join(binDir, "tapelock"), "#!/usr/bin/env bash\necho \"args: $*\"\n", 0o755)

	cassette := filepath.Join(t.TempDir(), "ok.jsonl")
	writeFile(t, cassette, "{}\n", 0o644)

	out, err := runCheckScript(t, binDir, cassette, "--status", "200")
	if err != nil {
		t.Fatalf("check.sh failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "--status 200") {
		t.Fatalf("check.sh did not pass through extra args:\n%s", out)
	}
}
