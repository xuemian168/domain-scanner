English | [中文](./README.zh.md)

# Domain Scanner

A powerful and flexible domain name availability checker written in Go. This tool helps you find available domain names based on various patterns and filters.

## Features

- Check domain availability using DNS and WHOIS lookups
- Generate domains based on different patterns:
  - Pure numbers (e.g., 123.li)
  - Pure letters (e.g., abc.li)
  - Alphanumeric combinations (e.g., a1b.li)
- Advanced filtering using regular expressions
- Save results to separate files for available and registered domains
- Configurable delay between queries to avoid rate limiting

## Installation

```bash
git clone https://github.com/xuemian168/domain-scanner.git
cd domain-scanner
go mod tidy
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
- `-h`: Show help information

### Examples

1. Check 3-letter .li domains:
```bash
go run main.go -l 3 -s .li -p D
```

2. Check 3-digit .li domains:
```bash
go run main.go -l 3 -s .li -p d
```

3. Check domains containing 'abc':
```bash
go run main.go -l 5 -s .li -p D -r '.*abc.*'
```

4. Check domains starting with 'a' and ending with 'z':
```bash
go run main.go -l 4 -s .li -p D -r '^a.*z$'
```

5. Check domains containing only vowels:
```bash
go run main.go -l 3 -s .li -p D -r '^[aeiou]+$'
```

## Output

The program generates two output files:
- `available_domains_[pattern]_[length]_[suffix].txt`
- `registered_domains_[pattern]_[length]_[suffix].txt`

## Regular Expression Examples

1. `^a.*` - Domains starting with 'a'
2. `.*z$` - Domains ending with 'z'
3. `^[0-9]+$` - Pure numeric domains
4. `^[a-z]+$` - Pure alphabetic domains
5. `^[a-z][0-9][a-z]$` - Letter-number-letter pattern
6. `.*(com|net|org)$` - Domains with specific suffixes
7. `^[a-z]{2}\d{2}$` - Two letters followed by two numbers

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the GNU AFFERO GENERAL PUBLIC License - see the LICENSE file for details. 