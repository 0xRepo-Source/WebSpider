# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of WebSpider
- Advanced rate limiting with token bucket algorithm
- Special rate limiting mode for burst-then-block servers
- Regex-based URL filtering (accept/reject patterns)
- Two-phase operation (discover then download)
- Directory structure preservation
- Concurrent request limiting with semaphore control
- Comprehensive logging and verbose mode
- HTTP error detection and backoff
- Configurable timeouts and user agents

### Features
- **Core Functionality**
  - Web directory crawling and discovery
  - Selective file downloading
  - URL list generation and management
  
- **Rate Limiting**
  - Standard requests-per-second limiting
  - Special burst-then-block mode for sensitive servers
  - Automatic block detection and recovery
  - Configurable request windows and block durations

- **Filtering and Control**
  - Regex patterns for URL acceptance/rejection
  - Configurable crawl depth limits
  - File type targeting
  - Directory exclusion patterns

- **Reliability**
  - Graceful error handling
  - Connection timeout management
  - Automatic retry mechanisms
  - Progress tracking and logging

## [1.0.0] - 2025-11-08

### Added
- Initial stable release
- Complete feature set as described above
- Comprehensive documentation
- Cross-platform support (Windows, Linux, macOS)
- MIT license

---

## Release Notes Template

When creating a new release, use this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Bug fixes

### Security
- Security fixes
```