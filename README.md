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

- **Multi-method Verification**: Checks domain availability using multiple methods:
  - DNS records (NS, A, MX)
  - WHOIS information
  - SSL certificate verification
- **Advanced Filtering**: Filter domains using regular expressions
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
- `-r string`: Regex filter for domain names
- `-delay int`: Delay between queries in milliseconds (default: 1000)
- `-workers int`: Number of concurrent workers (default: 10)
- `-show-registered`: Show registered domains in output (default: false)
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

4. Use regex filter for specific patterns:
```bash
go run main.go -l 3 -s .li -p D -r "^[a-z]{2}[0-9]$"
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

## Error Handling

The tool includes robust error handling:
- Automatic retry mechanism for WHOIS queries (3 attempts)
- Timeout settings for SSL certificate checks
- Graceful handling of network issues
- Detailed error reporting

## Contributing

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[![AGPL-3.0 License](https://img.shields.io/badge/License-AGPL--3.0-green.svg)](LICENSE)

This project is licensed under the AGPL-3.0 License - see the [LICENSE](LICENSE) file for details. 