English | [中文](./README.zh.md)

# Domain Scanner

[![Go Version](https://img.shields.io/badge/go-1.22-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL--3.0-green.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/xuemian168/domain-scanner.svg?style=social)](https://github.com/xuemian168/domain-scanner/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/xuemian168/domain-scanner.svg?style=social)](https://github.com/xuemian168/domain-scanner/network/members)
[![GitHub Issues](https://img.shields.io/github/issues/xuemian168/domain-scanner.svg)](https://github.com/xuemian168/domain-scanner/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/xuemian168/domain-scanner.svg)](https://github.com/xuemian168/domain-scanner/pulls)

A powerful domain name availability checker written in Go. This tool helps you find available domain names by checking multiple registration indicators and providing detailed verification results.

### Web Version: [zli.li](https://zli.li)

![Star History Chart](https://api.star-history.com/svg?repos=xuemian168/domain-scanner&type=Date)

## Features

- **Dictionary Input Support**: Generate domains from word lists for practical domain checking
  - Read dictionary files (one word per line) for word-based domain generation
  - Smart mode detection between dictionary and pattern-based generation
  - Regex filtering support for dictionary words
- **Multi-method Verification**: Checks domain availability using multiple methods:
  - DNS records (NS, A, MX)
  - WHOIS information
  - SSL certificate verification
- **Advanced Filtering**: Filter domains using powerful regular expressions with regexp2 support
  - Backreferences for patterns like repeating characters
  - Lookarounds and Unicode properties
  - Perl-compatible regex syntax
- **Security Enhanced**: Built-in protection against ReDoS attacks
- **Performance Warning System**: Intelligent warnings for large domain scans with detailed impact analysis
- **Smart Scan Estimation**: Automatic calculation of scan time, network load, and resource usage
- **User Safety Protection**: Prevents accidental execution of multi-day scan operations
- **Concurrent Processing**: Multi-threaded domain checking with configurable worker count
- **Smart Error Handling**: Automatic retry mechanism for failed queries
- **Detailed Results**: Shows verification signatures for registered domains
- **Progress Tracking**: Real-time progress display with current/total count
- **File Output**: Saves results to separate files for available and registered domains
- **Configurable Delay**: Adjustable delay between queries to prevent rate limiting

## Installation

```bash
git clone https://github.com/xuemian168/domain-scanner.git
cd domain-scanner
go mod download
```

## Usage

```bash
go run main.go [options]
```

### Options

- `-l int`: Domain length (default: 3)
- `-s string`: Domain suffix (default: .li)
- `-p string`: Domain pattern:
  - `d`: Pure numbers (e.g., 123.li)
  - `D`: Pure letters (e.g., abc.li)
  - `a`: Alphanumeric (e.g., a1b.li)
- `-r string`: Regex filter for domain name prefix (supports advanced regexp2 features)
- `-dict string`: Dictionary file path (one word per line) for word-based domain generation
- `-delay int`: Delay between queries in milliseconds (default: 1000)
- `-workers int`: Number of concurrent workers (default: 10)
- `-show-registered`: Show registered domains in output (default: false)
- `-force`: Skip performance warnings for large domain sets (default: false)
- `-h`: Show help information

### Examples

1. Check 3-letter .li domains with 20 workers:
```bash
go run main.go -l 3 -s .li -p D -workers 20
```

2. Check domains with custom delay and workers:
```bash
go run main.go -l 3 -s .li -p D -delay 500 -workers 15
```

3. Show both available and registered domains:
```bash
go run main.go -l 3 -s .li -p D -show-registered
```

4. Use regex filter to match domain prefix:
```bash
go run main.go -l 3 -s .li -p D -r "^[a-z]{2}[0-9]$"
```

5. Find domains starting with specific letters:
```bash
go run main.go -l 5 -s .li -p D -r "^abc"
```

6. Use dictionary file for word-based domain checking:
```bash
go run main.go -dict words.txt -s .com
```

7. Use dictionary with regex filter for precise matching:
```bash
go run main.go -dict words.txt -s .com -r "^[a-z]{4,8}$"
```

8. Use advanced regexp2 features (backreferences for repeating patterns):
```bash
# Find domains with pattern like "aaa", "bbb", "ccc" (same letter repeated)
go run main.go -l 3 -s .li -p D -r "^(.)\1{2}$"

# Find domains with pattern like "abab", "cdcd" (two letters repeated)
go run main.go -l 4 -s .li -p D -r "^(..)\1$"
```

9. Skip performance warning for large domain sets:
```bash
go run main.go -l 7 -s .li -p D -force
```

## Performance Warning System

The tool includes an intelligent performance warning system to protect users from accidentally running extremely large scans:

### When Warnings Are Triggered
- Automatically triggered when domain length (`-l`) is greater than 5
- Displays detailed impact analysis before starting the scan

### Warning Information Provided
- **Domain Count**: Exact number of domains to be scanned
- **Time Estimation**: Calculated scan duration based on your settings
- **Network Load**: Total number of network requests that will be made
- **Resource Impact**: Memory and CPU usage warnings

### Example Warning Output
```
⚠️  PERFORMANCE WARNING ⚠️
═══════════════════════════════════════════════════════
You are about to scan 308915776 domains with the following settings:
• Pattern: D (charset size: 26)
• Length: 6 characters
• Workers: 10
• Delay: 1000 ms between queries

📊 Estimated Impact:
• Scan time: ~357.0 days (8580.0 hours)
• Network requests: 308915776 total
• Memory usage: High (processing 308915776 domains)

💡 Recommendations:
• Use regex filter (-r) to narrow down the search
• Consider shorter domain length (-l)
• Increase workers (-workers) for faster processing
• Decrease delay (-delay) if your network can handle it
• Use -force flag to skip this warning next time
═══════════════════════════════════════════════════════

Do you want to continue? (y/N):
```

### Bypassing Warnings
Use the `-force` flag to skip performance warnings:
```bash
go run main.go -l 6 -s .com -p D -force
```

## Output Format

### Progress Display
```
[1/100] Domain abc.com is AVAILABLE!
[2/100] Domain xyz.com is REGISTERED [DNS_NS, WHOIS]
[3/100] Domain 123.com is REGISTERED [DNS_A, SSL]
```

### Verification Signatures
- `DNS_NS`: Domain has name server records
- `DNS_A`: Domain has IP address records
- `DNS_MX`: Domain has mail server records
- `WHOIS`: Domain is registered according to WHOIS
- `SSL`: Domain has a valid SSL certificate

### Output Files
- Available domains: `available_domains_[pattern]_[length]_[suffix].txt`
- Registered domains: `registered_domains_[pattern]_[length]_[suffix].txt`

## Advanced Regex Features

This tool uses the powerful `regexp2` library, providing advanced regex capabilities:

### Backreferences
Match previously captured groups using `\1`, `\2`, etc. All regex patterns match domain prefix only:
- `^(.)\1{2}$` - Matches domain prefixes like "aaa", "bbb" (same character repeated 3 times)
- `^(..)\1$` - Matches domain prefixes like "abab", "cdcd" (two characters repeated)
- `^(.)(..)\1\2$` - More complex backreference patterns

### Safety Features
- **ReDoS Protection**: Built-in timeout protection (100ms) prevents catastrophic backtracking
- **Input Validation**: Automatically rejects potentially dangerous regex patterns
- **Complexity Limits**: Maximum 200 characters, limited quantifiers
- **Error Handling**: Graceful handling of regex matching errors

### Security Guidelines
⚠️ **Important**: Be careful with complex regex patterns to avoid performance issues.

**Safe patterns:**
```bash
# Simple character classes and quantifiers
-r "^[a-z]{2}[0-9]$"
-r "^(test|demo|temp)"
-r "^[a-z]*[0-9]+$"
```

**Potentially dangerous patterns (automatically blocked):**
```bash
# These patterns are blocked for security
-r "(.*)*"     # Nested quantifiers
-r "(.+)+"     # Catastrophic backtracking
-r "(a+)+"     # ReDoS attack pattern
```

## Error Handling

The tool includes robust error handling:
- Automatic retry mechanism for WHOIS queries (3 attempts)
- Timeout settings for SSL certificate checks
- Regex timeout protection (100ms) against ReDoS attacks
- Input validation for regex patterns
- Graceful handling of network issues
- Detailed error reporting

## Contributing

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

We welcome contributions from the community! Whether you're fixing bugs, adding new features, improving documentation, or reporting issues, your help is appreciated.

### How to Contribute

1. **Fork the Repository**: Create your own copy of the project
2. **Create a Feature Branch**: Work on your changes in a dedicated branch
3. **Make Your Changes**: Follow our coding guidelines and test thoroughly
4. **Submit a Pull Request**: Describe your changes and link any related issues

For detailed contribution guidelines, development setup, and coding standards, please read our [CONTRIBUTING.md](./CONTRIBUTING.md) file.

### Quick Start for Contributors

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/domain-scanner.git
cd domain-scanner

# Set up development environment
go mod download
go build -o domain-scanner main.go

# Create a feature branch
git checkout -b feature/your-feature-name

# Make changes and test
go run main.go -l 2 -s .test -p D

# Commit and push
git commit -m "feat: your feature description"
git push origin feature/your-feature-name
```

## License

[![AGPL-3.0 License](https://img.shields.io/badge/License-AGPL--3.0-green.svg)](LICENSE)

This project is licensed under the AGPL-3.0 License - see the [LICENSE](LICENSE) file for details. 

## Recent Updates

### v1.3.4 - 2025-09-02
- **Dictionary Input**: New `-dict` parameter for word-based domain generation from text files
- **Smart Mode Detection**: Intelligent switching between dictionary and pattern-based generation
- **Enhanced Flexibility**: Support for regex filtering on dictionary words for precise domain matching
- **Practical Usage**: Enable checking of real-world domains using common word lists
- **Documentation**: Comprehensive examples and usage guidelines for dictionary mode

### v1.3.3 - 2025-09-02
- **Performance Warnings**: Intelligent warnings for large domain scans with detailed impact analysis
- **User Safety**: Protection against accidental multi-day scan operations with confirmation prompts
- **Windows Fix**: Resolved critical release binary execution issues causing empty results
- **Stability**: Fixed concurrent processing race conditions across all platforms

### v1.3.2 - 2025-08-26
- **Security**: Added ReDoS attack protection with regex timeout (100ms)
- **Security**: Implemented regex complexity validation and dangerous pattern detection
- **Performance**: Restored memory-efficient streaming architecture
- **Enhancement**: Upgraded to regexp2 for advanced regex features (backreferences, lookarounds)
- **Enhancement**: Added comprehensive regex safety guidelines and examples
- **Stability**: Improved error handling for regex matching operations

### v1.3.1 - 2025-08-24
- **Added**: Multiple WHOIS server support for improved reliability
- **Added**: Exponential backoff retry mechanism for WHOIS queries  
- **Added**: Comprehensive reserved domain indicators (139 patterns)
- **Performance**: Reduced false positive rate by 67% (15% → 5%)
- **Performance**: Improved WHOIS query success rate by 23% (~75% → ~92%)

### v1.3.0
- **Performance Optimizations**: Significantly improved domain checking speed
- **Bug Fixes**: Fixed WHOIS parsing for .de domains and other TLDs
- **Code Quality**: Refactored internal architecture for better maintainability

📋 **[View Complete Changelog](docs/CHANGELOG.md)** - See detailed version history, technical improvements, and all changes.
