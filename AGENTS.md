# synod

### Tooling

- Define project commands in `Taskfile.yml`; use `mise.toml` to pin tool versions.
- Run `task fmt` and `task lint` when writing code.
- Use `task test` when executing tests.

### Comments

- Do not add comments unless the behavior is HUGELY surprising and therefore needs an explanation.

### Documentation

- For local setup and validation, read [docs/development.md](docs/development.md).
- For release or deployment changes, read [docs/deployment.md](docs/deployment.md).
