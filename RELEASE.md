Release checklist
=================

Follow this checklist before creating a release tag. This helps keep releases consistent and safe.

Local checks
- Run all tests and linters:

```bash
make ci-check
```

- Build a local snapshot to verify artifacts:

```bash
make release
```

- Update `CHANGELOG.md` or release notes with user-facing changes.

Tagging and publishing
- Create a signed or annotated tag following semver, e.g. `v1.2.3`:

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

- The `Release` GitHub Actions workflow will run on push of `v*.*.*` tags. The release job uses the `release` environment and will require manual approval if your repository config requires it.

Repository admin steps (enable environment approvals)
- In the GitHub repository settings > Environments, create an environment named `release`.
- Configure protection rules and required reviewers for this environment (e.g., select specific users or teams who must approve publish jobs).

Post-release
- Verify the release artifacts on the GitHub release page and update distribution docs if needed.

Security note
- The `Release` workflow uses the repository `GITHUB_TOKEN` to publish. If you need finer control, consider using a dedicated deploy key or GitHub App with limited scope.
