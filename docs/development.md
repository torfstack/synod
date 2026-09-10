# Development

## Setup

Run `mise install` to install the tools pinned in `mise.toml`, then
`task frontend-install` to install frontend dependencies from the lockfile.
Ensure your shell uses the mise tools, or prefix commands with `mise exec --`.

## Running locally

Docker Compose mounts the repository's `config.yaml` into the application.
Configure it for the Compose database at `synod-db:5432` and HTTP port 4000.
The database is `synoddb`, with user `synod` and password `synodpass`.

```sh
task start-local
```

This builds and starts the application and PostgreSQL in the background.
The application is exposed on port 4000, the debugger on 4200, and PostgreSQL
on host port 4100.

```sh
task stop-local
```

## Checks

```sh
task frontend-install
task fmt
task lint
task test
```

Tests include version-selection fixtures and the Go suite. Database integration
tests require Docker. Lint includes the frontend and Helm chart.
