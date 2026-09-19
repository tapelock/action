# Tapelock Action

Run Tapelock in GitHub Actions.

The action installs Tapelock, runs your configured checks, and exposes the
result through the GitHub Actions workflow.

**Status:** v0.1.0 - early, expect breaking changes to inputs and behavior.

## Usage

```yaml
name: LLM Tests

on:
  pull_request:

jobs:
  tapelock:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: tapelock/action@v0.1.0
        with:
          cassette: cassettes/
```

The action fails the workflow when a Tapelock check fails.

## What it does

* Installs the Tapelock CLI
* Runs Tapelock checks in CI
* Preserves Tapelock exit codes
* Exposes failures directly in the workflow

For configuration and usage, see the [Tapelock documentation](https://github.com/tapelock/tapelock).

## Development

```bash
go test ./...
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
