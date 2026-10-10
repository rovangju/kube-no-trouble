# AGENTS.md

Guidance for AI coding agents working on kubent — a CLI that detects
deprecated Kubernetes APIs in clusters and manifests.

## Development Checks

- Run `make test` before proposing changes — it enforces formatting gates
  (code and rego rules) and runs the unit tests
- Rego rules must stay formatted via the Makefile's format targets

## Conventional Commits

Commits follow [Conventional Commits](https://www.conventionalcommits.org/):
`<type>: <summary>` — imperative, present tense, capitalized, no trailing
period. Common types: feat, fix, docs, ref, style, test, chore, ci, build,
dep. Breaking changes: `!` suffix and a `BREAKING CHANGE:` footer.
(No mechanical enforcement — follow the convention in good faith.)

## Releases (CalVer)

A release = pushing a tag in `YY.MM.BUILD` form (e.g. `26.10.1`; build
resets to 1 each January). Release notes are generated automatically by
GitHub from merged PRs. Nightly = a single rolling `nightly` tag, published
as a pre-release. No hotfixes — fixes ride the next build.

## Maintenance: adding deprecation rules for a new K8s version

The most common maintenance task. Worked example: PR #1.

Supported window: the latest official K8s minor + the 5 before it (policy
stated in the README "Supported Versions" section). When adding a new
version, prune rules, fixtures, tests, and matrix entries that fall below
the window.

1. **Ground truth = API schemas.** Read the type definitions in
   `kubernetes/kubernetes` at the target version's tag — e.g.
   `staging/src/k8s.io/api/<group>/<version>/zz_generated.prerelease-lifecycle.go`
   for API-level deprecations/removals, and `types.go` in the same package
   for field-level deprecations. Release notes, blog posts, and CHANGELOGs
   are hints only — confirm every finding against the schema.
2. **Rules:** add a rego file per version in `pkg/rules/rego/` and cite
   the source URL(s) in a header comment (existing convention).
3. **Fixtures + tests:** add example manifests under `fixtures/` and a
   per-version test under `test/` — fixture cases plus version-filtering
   assertions.
4. **Integration matrix:** keep the CI workflow's K8s version matrix
   current (add the new version, drop EOL ones). Bump
   `MinSupportedTargetVersion` in `pkg/config/config.go` to match the
   window — it drives the below-window warnings for target and cluster
   versions.
5. **Double-check pass** before merging: audit for overlooked
   deprecations and nuances; schemas win over prose.

## Do NOT

- Add container image builds/publishing (the image pipeline was removed)
- Tag new releases in semver form — releases are CalVer only
