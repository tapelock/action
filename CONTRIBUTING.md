# Contributing to Tapelock Action

This action is in early development (pre-v1); expect breaking changes to
its inputs and behavior.

## Scope for v1

The core path this action must get right, in order:

```
GitHub Action
     ↓
tapelock check
     ↓
exit code
     ↓
GitHub ✓ / ✗
```

Installing the pinned `tapelock` CLI, running it against the configured
cassette, and turning its exit code into a passing or failing job is the
whole job of v1. A PR comment summarizing the result is a UX improvement
on top of that, not part of the core path — do not document or promise it
until it is actually implemented.

## Development

```bash
go test ./...
```

## Reporting issues

Open a GitHub issue with the workflow YAML that triggers the problem and
the job's log output.
