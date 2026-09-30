# Releasing

The release workflow (`.github/workflows/release.yml`) runs `goreleaser` on a `v*` tag and publishes the binaries. The `curl | sh` installer reads the latest release.

## Steps

1. Make sure `main` is green (`cli` and `gates`).
2. Update `CHANGELOG.md`: move the `Unreleased` entries under the new version and date.
3. Pick the version (see below).
4. Land the changelog through a pull request.
5. Tag and push:

   ```sh
   git tag -a v1.0.0 -m "Release 1.0.0"
   git push origin v1.0.0
   ```

6. Watch the release workflow, then confirm the install:

   ```sh
   curl -fsSL https://raw.githubusercontent.com/Kerman-Sanjuan/octospec/main/install.sh | sh
   octospec version
   ```

## Versioning

`vMAJOR.MINOR.PATCH`. A breaking change bumps the major, a new capability bumps the minor, a fix bumps the patch. The tag is the source of truth, so the version is not hand-edited in scattered files.
