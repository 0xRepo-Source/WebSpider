# WebSpider

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xRepo-Source/WebSpider)](https://goreportcard.com/report/github.com/0xRepo-Source/WebSpider)
[![Release](https://img.shields.io/github/v/release/0xRepo-Source/WebSpider)](https://github.com/0xRepo-Source/WebSpider/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/0xRepo-Source/WebSpider/build.yml)](https://github.com/0xRepo-Source/WebSpider/actions)

A sophisticated web directory crawler built in Go that provides advanced rate limiting, intelligent filtering, and selective downloading capabilities. Designed as a modern alternative to `wget --spider` with enhanced features for respectful web crawling.

## Features

- **Advanced Rate Limiting**: Supports both standard requests-per-second limiting and specialized burst-with-backoff patterns
- **Two-Phase Operation**: Discover directory structure first, then selectively download only desired files
- **Intelligent Filtering**: Powerful regex-based URL acceptance and rejection patterns
- **Respectful Crawling**: Built-in detection of rate limiting responses with automatic backoff
- **Directory Structure Preservation**: Maintains original directory hierarchy during downloads
- **Concurrent Processing**: Configurable concurrent request limiting with semaphore control
- **Comprehensive Logging**: Detailed progress tracking and verbose debugging options

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/0xRepo-Source/WebSpider.git
cd WebSpider

# Build from source
go mod tidy
go build -o webspider

# Or install directly
go install github.com/0xRepo-Source/WebSpider@latest
```

### Basic Usage

```bash
# Discover directory structure (recommended first step)
./webspider -url "https://example.com/files/" -discover-only -verbose

# Edit the generated discovered-urls.txt file to select desired files

# Download selected files
./webspider -urls "discovered-urls.txt" -rate 0.5
```

## Advanced Usage

### Standard Rate Limiting

For servers with standard rate limiting policies:

For servers with standard rate limiting policies:

```bash
# Conservative crawling
./webspider -url "https://example.com/docs/" -rate 0.5 -discover-only

# Moderate speed crawling
./webspider -url "https://example.com/files/" -rate 2.0 -depth 4
```

### Special Rate Limiting

For servers that implement burst-then-block rate limiting (e.g., 2 requests per 5 seconds, then 10-second block):

```bash
# Default special rate limiting (2 req/5sec, 10sec block)
./webspider -url "https://sensitive-server.com/" -special-rate -discover-only

# Custom burst limiting (3 req/10sec, 15sec block)
./webspider -url "https://custom-server.com/" \
  -special-rate \
  -max-requests 3 \
  -time-window 10s \
  -block-duration 15s \
  -verbose
```

### Content Filtering

Target specific file types and exclude unwanted content:

```bash
# Academic papers and documents
./webspider -url "https://university.edu/publications/" \
  -accept "\.(pdf|doc|docx|ppt|pptx)$" \
  -discover-only \
  -save-list "academic-papers.txt"

# Software packages only
./webspider -url "https://releases.example.com/" \
  -accept "\.(tar\.gz|zip|deb|rpm|dmg)$" \
  -reject "/archive/|/old/" \
  -discover-only

# Exclude web assets
./webspider -url "https://docs.example.com/" \
  -reject "\.(css|js|jpg|jpeg|png|gif|svg|ico)$" \
  -depth 5
```

## Command Line Reference

### Core Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-url` | Base URL to crawl | Required | `https://example.com/files/` |
| `-urls` | File containing URLs to download | - | `discovered-urls.txt` |
| `-discover-only` | Only discover URLs, don't download | `false` | - |
| `-depth` | Maximum crawling depth | `3` | `5` |
| `-output` | Output directory for downloads | `./downloads` | `./my-files` |
| `-save-list` | File to save discovered URLs | `discovered-urls.txt` | `results.txt` |
| `-verbose` | Enable verbose logging | `false` | - |

### Rate Limiting Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-rate` | Standard rate limit (req/sec) | `1.0` | `0.5` |
| `-special-rate` | Enable burst-then-block limiting | `false` | - |
| `-max-requests` | Max requests in time window | `2` | `3` |
| `-time-window` | Time window for request limiting | `5s` | `10s` |
| `-block-duration` | Duration server blocks after limit | `10s` | `15s` |

### Filtering Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-accept` | Regex pattern for URLs to accept | - | `\.(pdf\|doc)$` |
| `-reject` | Regex pattern for URLs to reject | - | `\.(css\|js)$` |
| `-user-agent` | Custom User-Agent string | `WebSpider/1.0` | `Mozilla/5.0...` |
| `-timeout` | HTTP request timeout | `30s` | `60s` |

## Use Cases

### Academic Research
Download research papers and documentation while respecting university server policies:

```bash
./webspider -url "https://university.edu/papers/" \
  -special-rate \
  -accept "\.(pdf|doc|docx)$" \
  -discover-only \
  -verbose
```

### Software Distribution
Mirror software releases with conservative rate limiting:

```bash
./webspider -url "https://releases.project.org/" \
  -rate 0.5 \
  -accept "\.(tar\.gz|zip|deb|rpm)$" \
  -depth 3 \
  -output "./software-mirror"
```

### Documentation Archival
Archive website documentation excluding assets:

```bash
./webspider -url "https://docs.example.com/" \
  -reject "\.(jpg|jpeg|png|gif|css|js)$" \
  -accept "\.(html|htm|pdf|txt|md)$" \
  -rate 2.0
```

## Best Practices

### Respectful Crawling
- Always start with `-discover-only` to understand the site structure
- Use conservative rate limits (`0.5-1.0` req/sec) for unknown servers
- Monitor server responses with `-verbose` flag
- Respect robots.txt when possible

### Efficient Filtering
- Use `-accept` patterns to target specific file types early
- Combine with `-reject` patterns to exclude unwanted content
- Set appropriate `-depth` limits to avoid unnecessary crawling
- Test regex patterns before large crawls

### Error Recovery
- Use `-timeout` for unreliable connections
- Enable `-verbose` logging for debugging
- Check generated URL lists before downloading
- Consider using `-special-rate` for sensitive servers

## Technical Details

### Rate Limiting Implementation
- **Standard Mode**: Token bucket algorithm via `golang.org/x/time/rate`
- **Special Mode**: Sliding window with automatic block detection
- **Backoff Strategy**: Exponential backoff on HTTP 429/503 responses
- **Concurrent Control**: Configurable semaphore limiting

### URL Discovery
- **HTML Parsing**: Uses `goquery` for robust DOM traversal
- **Link Resolution**: Handles relative and absolute URLs correctly
- **Deduplication**: Prevents revisiting discovered URLs
- **Filtering**: Real-time regex-based URL filtering

### Download Management
- **Directory Preservation**: Maintains original site structure
- **Atomic Operations**: Safe concurrent file creation
- **Error Handling**: Graceful failure recovery
- **Progress Tracking**: Comprehensive logging system

## Troubleshooting

### Common Issues

**Rate Limiting Errors (HTTP 429)**
```bash
# Solution: Use special rate limiting mode
./webspider -special-rate -max-requests 1 -time-window 10s
```

**Connection Timeouts**
```bash
# Solution: Increase timeout and reduce rate
./webspider -timeout 60s -rate 0.1
```

**Memory Usage with Large Sites**
```bash
# Solution: Limit depth and use filtering
./webspider -depth 2 -accept "\.(pdf|zip)$"
```

### Getting Help

- Use `-h` flag to see all available options
- Enable `-verbose` logging for detailed operation information
- Check the generated URL list files for unexpected results
- Test with small depth values before full crawls

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
git clone https://github.com/0xRepo-Source/WebSpider.git
cd WebSpider
go mod tidy
go build -o webspider
```

### Running Tests

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing
- Rate limiting powered by [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)
- Inspired by wget's spider functionality with modern improvements