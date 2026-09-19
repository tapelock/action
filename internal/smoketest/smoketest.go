// Package smoketest exercises scripts/check.sh, the action's core path
// (install tapelock, run tapelock check, let its exit code decide the
// job), against a fake `tapelock` binary, without needing a real
// tapelock release or a GitHub Actions runner.
package smoketest
