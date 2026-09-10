# Release and deployment

Run `mise install` to install the pinned tools. Docker with Buildx is also required
for publishing images. Authenticate Docker to `ghcr.io` before publishing locally.

Pushes to `main` run lint and tests, calculate the next patch version once, publish
that image, and tag the exact commit. Releases are serialized; GitHub may replace
an older pending run with a newer one, so not every intermediate commit is
necessarily released. An already tagged commit is rejected on rerun.

`task next-version` considers only exact `vMAJOR.MINOR.PATCH` tags, chooses the
highest numeric version, and increments its patch. With no matching tags it returns
`v0.0.1`. It reads local tags; CI refreshes them before allocating a version.

To publish manually, choose a fresh version explicitly:

```sh
task release VERSION=v0.0.55
```

This task only publishes an image; CI owns release Git tags. Coordinate manual
publishes with CI and never reuse a published version. Images default to
`linux/amd64`, matching the existing GitHub runner builds; override `PLATFORM`
when targeting another architecture.

Deploy an existing published image using an explicit cluster and namespace:

```sh
task deploy VERSION=v0.0.55 KUBE_CONTEXT=my-cluster NAMESPACE=default
```

The namespace must already exist and contain the `synod` configuration secret and
`regcred` image pull secret. Helm stores the release in that same namespace.
The deployment still expects the application's configured HTTP port to be 4000.

Deploy local changes, including uncommitted changes, with a unique image tag:

```sh
task dev-deploy KUBE_CONTEXT=my-cluster NAMESPACE=default
```

The generated tag contains the commit, UTC timestamp and random suffix. Both tasks
receive the same tag, so each invocation changes the Pod template and starts a
rollout. This targets the same `synod` release as `task deploy`; use a separate
namespace for a separate development instance. The existing fixed `hostPort: 4000`
also requires those instances to run on different nodes.

Helm waits for readiness and rolls back failed upgrades, with a default timeout of
five minutes (`DEPLOY_TIMEOUT=10m` overrides it). `/readyz` becomes available after
startup migrations and initialization complete; it verifies HTTP responsiveness,
not ongoing database connectivity or whether the vault is unsealed. The deployment
uses `Recreate`, so upgrades still involve downtime. Older images without `/readyz`
require their matching older chart; Helm rollback restores the previous chart
as well as its image. Helm rollback does not undo
database migrations; migrations must remain compatible with the previous app.

If publishing succeeds but pushing the Git tag fails, verify that the published
image belongs to that workflow commit and push that exact version tag to that
commit before rerunning CI. Image publication and Git tagging are not atomic.

