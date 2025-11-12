# CI/CD Pipeline

This project uses GitHub Actions and GoReleaser for automated testing, building, and releasing.

## Overview

The CI/CD pipeline consists of three main workflows:

1. **Test Workflow** - Runs on every push and PR
2. **Release Workflow** - Runs when a version tag is pushed
3. **Nightly Build Workflow** - Runs automatically on a schedule

## Workflows

### Test Workflow (`.github/workflows/test.yaml`)

Runs automatically on:
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**
- **test**: Runs tests across multiple platforms (Ubuntu, macOS, Windows) with Go 1.25
  - Downloads dependencies
  - Runs `go vet`
  - Runs tests with race detector and coverage
  - Uploads coverage to Codecov (Ubuntu only)

- **lint**: Runs golangci-lint for code quality checks

- **build**: Verifies the project builds successfully

- **goreleaser-check**: Validates the GoReleaser configuration

### Release Workflow (`.github/workflows/release.yaml`)

Runs automatically when a version tag (e.g., `v1.0.0`) is pushed.

**What it does:**
1. Checks out code with full git history
2. Sets up Go 1.25
3. Sets up ko for container image building
4. Logs into GitHub Container Registry (GHCR)
5. Runs GoReleaser to:
   - Build binaries for multiple platforms (Linux, macOS, Windows) and architectures (amd64, arm64)
   - Create universal binaries for macOS
   - Generate archives (tar.gz for Unix, zip for Windows)
   - Create container images using ko (multi-arch: linux/amd64, linux/arm64)
   - Publish to Homebrew tap (if configured)
   - Generate changelog
   - Create GitHub Release with all artifacts
   - Generate checksums
6. Attests build provenance for supply chain security

**Artifacts produced:**
- Binary archives for Linux, macOS, Windows (multiple architectures)
- Container images: `ghcr.io/spandigital/token-visualizer:TAG` and `:latest`
- Homebrew formula (when configured)
- Checksums and signatures
- SBOM (Software Bill of Materials)

### Nightly Build Workflow (`.github/workflows/nightly.yaml`)

Runs automatically:
- Tuesday through Saturday at midnight UTC
- Can be triggered manually

**What it does:**
1. Creates/updates a `nightly` tag
2. Builds snapshot releases (not published to GitHub Releases)
3. Uploads artifacts for 7 days
4. Cleans up old nightly releases (keeps latest 5)

## GoReleaser Configuration

The `.goreleaser.yaml` file defines the release process:

### Builds
- **Platforms**: Linux, macOS, Windows
- **Architectures**: amd64, arm64, arm (v7)
- **Flags**: `-trimpath` for reproducible builds
- **LDflags**: Embeds version, commit, and build date

### Universal Binaries
Creates universal binaries for macOS (Intel + ARM combined)

### Archives
- Format: `.tar.gz` for Unix, `.zip` for Windows
- Includes: README.md, LICENSE, CLAUDE.md

### Container Images (ko)
Uses [ko](https://ko.build/) for building container images:
- **Base**: Minimal distroless base images
- **Multi-arch**: linux/amd64, linux/arm64
- **Registry**: GitHub Container Registry (ghcr.io)
- **Tags**: Version tag and `latest` (for non-prerelease)
- **SBOM**: Generates SPDX SBOM for supply chain security
- **Labels**: OCI standard labels for metadata

### Homebrew Tap
Publishes to `SPANDigital/homebrew-tap` (requires `GH_PAT` secret)

### Changelog
Automatically generated from commit messages, grouped by:
- Features (commits starting with `feat:`)
- Bug fixes (commits starting with `fix:`)
- Enhancements (commits starting with `enhance:`)

## Secrets Required

### Repository Secrets

- **GITHUB_TOKEN**: Automatically provided by GitHub Actions
- **GH_PAT** (optional): Personal Access Token with `repo` and `write:packages` scope
  - Required for Homebrew tap publishing
  - Can be created at: https://github.com/settings/tokens

- **CODECOV_TOKEN** (optional): For uploading coverage reports
  - Sign up at: https://codecov.io/
  - Add your repository and get the token

## Making a Release

### 1. Update Version

Ensure your code is ready for release and committed to `main`.

### 2. Create and Push a Tag

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag to trigger the release workflow
git push origin v1.0.0
```

### 3. Monitor the Release

1. Go to the **Actions** tab in GitHub
2. Watch the **Release** workflow run
3. Once complete, check the **Releases** page for the new release

### 4. Verify Artifacts

The release should include:
- Binary archives for all platforms
- Container images on GHCR
- Homebrew formula (if configured)
- Changelog
- Checksums

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0 → v2.0.0): Incompatible API changes
- **MINOR** version (v1.0.0 → v1.1.0): New functionality (backwards compatible)
- **PATCH** version (v1.0.0 → v1.0.1): Bug fixes (backwards compatible)

## Container Usage

Once published, the container images can be used:

```bash
# Pull the latest image
docker pull ghcr.io/spandigital/token-visualizer:latest

# Run it
echo "Hello, world!" | docker run -i ghcr.io/spandigital/token-visualizer:latest

# Specific version
docker pull ghcr.io/spandigital/token-visualizer:v1.0.0
```

## Local Testing

Test the release process locally (without publishing):

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Test build (doesn't publish)
export KO_DOCKER_REPO=ghcr.io/spandigital/token-visualizer
goreleaser build --snapshot --clean

# Full release dry-run
goreleaser release --snapshot --clean --skip=publish
```

## Troubleshooting

### GoReleaser Check Fails

```bash
# Validate configuration
goreleaser check

# Build a snapshot to test
export KO_DOCKER_REPO=ghcr.io/spandigital/token-visualizer
goreleaser build --snapshot --clean
```

### Tests Fail

```bash
# Run tests locally
go test ./...

# Run with race detector
go test -race ./...

# Run linter
golangci-lint run
```

### Container Build Fails

Ensure ko is installed:
```bash
go install github.com/google/ko@latest
```

## Maintenance

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update GitHub Actions
# Check .github/workflows/*.yaml for version updates
```

### Updating GoReleaser Config

When updating `.goreleaser.yaml`:
1. Test locally: `goreleaser check`
2. Run a snapshot build: `goreleaser build --snapshot --clean`
3. Commit and push
4. Test with a pre-release tag (e.g., `v1.0.0-beta.1`)

## Best Practices

1. **Always test locally** before pushing tags
2. **Use conventional commits** for better changelogs
3. **Keep dependencies updated** regularly
4. **Test on multiple platforms** before releasing
5. **Document breaking changes** in commit messages
6. **Use pre-release versions** (v1.0.0-beta.1) for testing

## References

- [GoReleaser Documentation](https://goreleaser.com/)
- [ko Documentation](https://ko.build/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Conventional Commits](https://www.conventionalcommits.org/)
