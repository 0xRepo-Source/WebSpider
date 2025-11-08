# WebSpider

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xRepo-Source/WebSpider)](https://goreportcard.com/report/github.com/0xRepo-Source/WebSpider)
[![Release](https://img.shields.io/github/v/release/0xRepo-Source/WebSpider)](https://github.com/0xRepo-Source/WebSpider/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/0xRepo-Source/WebSpider/build.yml)](https://github.com/0xRepo-Source/WebSpider/actions)
[![Downloads](https://img.shields.io/github/downloads/0xRepo-Source/WebSpider/total?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider/releases)
[![GitHub stars](https://img.shields.io/github/stars/0xRepo-Source/WebSpider?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/0xRepo-Source/WebSpider?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider/network/members)
[![GitHub issues](https://img.shields.io/github/issues/0xRepo-Source/WebSpider?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider/issues)
[![GitHub last commit](https://img.shields.io/github/last-commit/0xRepo-Source/WebSpider?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider/commits/main)
[![Code size](https://img.shields.io/github/languages/code-size/0xRepo-Source/WebSpider?style=flat&logo=github)](https://github.com/0xRepo-Source/WebSpider)
[![Top language](https://img.shields.io/github/languages/top/0xRepo-Source/WebSpider?style=flat&logo=go)](https://github.com/0xRepo-Source/WebSpider)
[![Robots.txt Compliant](https://img.shields.io/badge/robots.txt-compliant-green?style=flat&logo=robot)](https://www.robotstxt.org/)
[![Cross Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat)](https://github.com/0xRepo-Source/WebSpider/releases)

A sophisticated web directory crawler built in Go that provides advanced rate limiting, intelligent filtering, and selective downloading capabilities. Designed as a modern alternative to `wget --spider` with enhanced features for respectful web crawling.

## Features

- **Advanced Rate Limiting**: Supports both standard requests-per-second limiting and specialized burst-with-backoff patterns
- **Robots.txt Support**: Automatic robots.txt fetching, parsing, and compliance with crawl-delay and disallow/allow rules
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

### Windows Setup for Easy Access

For convenient command-line usage on Windows, you can set up WebSpider to be accessible from anywhere:

1. **Download/Build WebSpider**:
   - Download `webspider-windows-amd64.exe` from the [releases page](https://github.com/0xRepo-Source/WebSpider/releases)
   - Or build from source as shown above

2. **Create Directory and Rename**:
   ```cmd
   # Create a dedicated directory
   mkdir C:\WebSpider
   
   # Copy the executable and rename it for easier typing
   copy webspider-windows-amd64.exe C:\WebSpider\ws.exe
   ```

3. **Add to PATH Environment Variable**:
   - Press `Win + R`, type `sysdm.cpl`, press Enter
   - Click "Environment Variables" button
   - Under "User variables" or "System variables", find and select "Path"
   - Click "Edit" → "New" → Add `C:\WebSpider`
   - Click "OK" to close all dialogs

4. **Usage from Anywhere**:
   ```cmd
   # Now you can use 'ws' from any directory
   ws -url "https://example.com/files/" -discover-only -verbose
   ws -urls "discovered-urls.txt" -rate 0.5
   ws -special-rate -url "https://sensitive-server.com/" -discover-only
   ```

### PowerShell Alternative Setup

If you prefer using PowerShell profiles:

```powershell
# Create a PowerShell function (add to your PowerShell profile)
function ws { & "C:\WebSpider\ws.exe" $args }

# Usage
ws -url "https://example.com/" -discover-only
```

### Basic Usage

```bash
# Discover directory structure (recommended first step)
./webspider -url "https://example.com/files/" -discover-only -verbose

# Edit the generated discovered-urls.txt file to select desired files

# Download selected files
./webspider -urls "discovered-urls.txt" -rate 0.5
```

**Windows (after PATH setup):**
```cmd
# Same commands but using the shorter 'ws' alias
ws -url "https://example.com/files/" -discover-only -verbose
ws -urls "discovered-urls.txt" -rate 0.5
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

### Robots.txt Compliance

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-ignore-robots` | Ignore robots.txt rules | `false` (respects robots.txt) | `-ignore-robots` |

**Robots.txt Features:**
- **Automatic fetching**: Downloads and parses robots.txt from each domain
- **User-agent matching**: Respects rules for `WebSpider`, `*`, and custom user agents
- **Crawl-delay support**: Automatically adjusts rate limiting based on `Crawl-delay` directive
- **Path filtering**: Honors `Disallow` and `Allow` path patterns
- **Caching**: Caches robots.txt for 24 hours to reduce server load

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

**Windows:**
```cmd
ws -url "https://university.edu/papers/" -special-rate -accept "\.(pdf|doc|docx)$" -discover-only -verbose
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

**Windows:**
```cmd
ws -url "https://releases.project.org/" -rate 0.5 -accept "\.(tar\.gz|zip|deb|rpm)$" -depth 3 -output ".\software-mirror"
```

### Documentation Archival
Archive website documentation excluding assets:

```bash
./webspider -url "https://docs.example.com/" \
  -reject "\.(jpg|jpeg|png|gif|css|js)$" \
  -accept "\.(html|htm|pdf|txt|md)$" \
  -rate 2.0
```

**Windows:**
```cmd
ws -url "https://docs.example.com/" -reject "\.(jpg|jpeg|png|gif|css|js)$" -accept "\.(html|htm|pdf|txt|md)$" -rate 2.0
```

### Robots.txt Compliant Crawling
Crawl a website while automatically respecting robots.txt rules:

```bash
./webspider -url "https://example.com/data/" \
  -verbose \
  -discover-only \
  -accept "\.(csv|json|xml)$"
```

**Windows:**
```cmd
ws -url "https://example.com/data/" -verbose -discover-only -accept "\.(csv|json|xml)$"
```

This will:
- Automatically fetch and parse robots.txt from example.com
- Respect any `Disallow` paths that block access to `/data/` or subdirectories
- Honor `Crawl-delay` directives by adjusting the rate limiter
- Skip URLs blocked by robots.txt (shown in verbose output)

To override robots.txt protection (use responsibly):
```bash
./webspider -url "https://example.com/data/" -ignore-robots -rate 0.5
```

## Best Practices

### Respectful Crawling
- Always start with `-discover-only` to understand the site structure
- WebSpider automatically respects robots.txt by default
- Use conservative rate limits (`0.5-1.0` req/sec) for unknown servers
- Monitor server responses with `-verbose` flag
- Only use `-ignore-robots` when you have permission to bypass robots.txt

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

### Windows-Specific Issues

**"ws is not recognized" Error**
```cmd
# Check if PATH was added correctly
echo %PATH%

# Verify the executable exists
dir C:\WebSpider\ws.exe

# Try full path if PATH isn't working
C:\WebSpider\ws.exe -url "https://example.com/" -discover-only
```

**Permission Issues on Windows**
```cmd
# Run Command Prompt as Administrator when setting up PATH
# Or use PowerShell profile method instead
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