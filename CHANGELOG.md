# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### 🐛 Fixed
- **Special rate limiting burst behavior**: Fixed issue where special rate limiting allowed 4-70 request bursts instead of controlled spacing
- **Sequential processing**: Special rate limiting now forces sequential requests (concurrentLimit=1) to prevent simultaneous bypassing
- **Minimum interval enforcement**: Now calculates and enforces minimum intervals between ALL requests (TimeWindow ÷ MaxRequests)
- **Improved timing accuracy**: Better request spacing prevents server blocking with consistent timing

### 🔧 Improved
- **Rate limiting logic**: Enhanced `waitForRateLimit` function with proper interval calculation
- **Verbose logging**: Added detailed logging for minimum interval enforcement
- **Concurrent control**: Dynamic concurrent limits based on rate limiting mode (1 for special, 5 for normal)

## [1.1.0] - 2025-11-08

### 🤖 Added - Robots.txt Support
- **Automatic robots.txt fetching**: Downloads and parses robots.txt from each target domain
- **User-agent matching**: Respects rules for WebSpider, wildcard (*), and custom user agents  
- **Crawl-delay support**: Automatically adjusts rate limiting based on robots.txt Crawl-delay directives
- **Path filtering**: Honors Disallow and Allow patterns during URL discovery and downloads
- **Domain caching**: 24-hour caching of robots.txt to reduce server load and improve performance
- **Override capability**: New `-ignore-robots` flag for cases where you have permission to bypass robots.txt
- **Verbose logging**: Detailed output showing robots.txt loading and rule applications
- **Graceful fallback**: Assumes all allowed when robots.txt is unavailable or inaccessible

### 📚 Documentation  
- Updated README.md with comprehensive robots.txt documentation
- Added robots.txt compliance examples and use cases
- Documented new `-ignore-robots` command-line flag
- Updated best practices section to emphasize robots.txt respect
- Added technical details about robots.txt implementation

### 🔧 Technical Improvements
- New `RobotsData` struct for structured robots.txt rule storage
- Robust robots.txt parser handling all standard directives
- Integration into both discovery and download phases
- Dynamic rate limiter adjustment based on crawl-delay values
- Domain-based caching system with TTL management
- Enhanced error handling for robots.txt fetch failures

### ✅ Testing
- Verified on real websites (GitHub.com with 59 robots.txt rules)
- Tested graceful handling of missing robots.txt files
- Confirmed proper rate limiting and URL filtering behavior
- Validated verbose output provides clear robots.txt information

### Breaking Changes
- None. All changes are backward compatible.
- Default behavior now includes robots.txt compliance (more restrictive but ethical)

## [1.0.0] - 2025-11-07

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

## [1.0.0] - 2025-11-07

### 🎉 Added - Initial Release
- Advanced rate limiting with token bucket algorithm
- Special rate limiting mode for burst-then-block servers
- Regex-based URL filtering (accept/reject patterns)
- Two-phase operation (discover then download)
- Directory structure preservation  
- Concurrent request limiting with semaphore control
- Comprehensive logging and verbose mode
- HTTP error detection and backoff
- Configurable timeouts and user agents
- Cross-platform builds (Windows, Linux, macOS)
- Professional CI/CD with GitHub Actions
- Windows PATH integration documentation

### Platform Support
- Windows (AMD64, ARM64)
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Automated builds and releases via GitHub Actions

---
