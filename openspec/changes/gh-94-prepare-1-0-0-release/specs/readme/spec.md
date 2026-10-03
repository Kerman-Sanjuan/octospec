## ADDED Requirements

### Requirement: The README displays a logo
The README SHALL display a logo, and the repository SHALL carry the logo asset as both SVG and PNG.

#### Scenario: A visitor sees the logo
- **WHEN** a visitor opens the README
- **THEN** the logo renders at the top, and both the SVG and PNG assets exist in the repository.

### Requirement: The README shows status badges
The README SHALL show status badges for CI, the latest release, the license, and the Go version.

#### Scenario: Badges are present
- **WHEN** a visitor opens the README
- **THEN** badges for CI, the latest release, the license, and the Go version appear near the top and resolve to live status.

### Requirement: The README links the docs, the changelog, and About
The README SHALL link the docs, `CHANGELOG.md`, and a place to learn about the project.

#### Scenario: Links are reachable
- **WHEN** a reader follows the README links
- **THEN** they reach the docs, the changelog, and the About page.
