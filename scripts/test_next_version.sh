#!/bin/sh
set -eu

script=$(pwd)/next_version
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' EXIT HUP INT TERM
cd "$fixture"

if sh "$script" >/dev/null 2>&1; then
    echo "next_version should fail outside a Git repository" >&2
    exit 1
fi
git init -q
# Keep fixtures independent of developer hooks and signing configuration.
git config core.hooksPath /dev/null
git config commit.gpgsign false
git config tag.gpgsign false
git -c user.name=Test -c user.email=test@example.com commit -q --allow-empty -m initial
expect_version() {
    actual=$(sh "$script")
    if [ "$actual" != "$1" ]; then
        printf 'expected %s, got %s\n' "$1" "$actual" >&2
        exit 1
    fi
}
expect_version v0.0.1
for tag in backend/v9.0.0 frontend/v8.0.0 v9.0.0-rc.1 latest v01.0.0; do
    git tag "$tag"
done
expect_version v0.0.1
for tag in v0.0.9 v0.0.10 v0.0.2; do
    git tag "$tag"
done
expect_version v0.0.11
git tag v0.10.0
git tag v0.9.99
expect_version v0.10.1
git tag v2.0.0
git tag v1.99.99
git -c user.name=Test -c user.email=test@example.com tag -a v0.1.0 -m 'newer tag, older version'
expect_version v2.0.1
printf 'next_version checks passed\n'
