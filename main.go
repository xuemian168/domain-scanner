package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/likexian/whois"
)

func checkDomainAvailability(domain string) (bool, error) {
	// First check DNS records
	// Try to resolve NS records
	nsRecords, err := net.LookupNS(domain)
	if err == nil && len(nsRecords) > 0 {
		return false, nil
	}

	// Try to resolve A records
	ipRecords, err := net.LookupIP(domain)
	if err == nil && len(ipRecords) > 0 {
		return false, nil
	}

	// If no DNS records found, check WHOIS as a secondary check
	result, err := whois.Whois(domain)
	if err != nil {
		return false, err
	}

	// Convert WHOIS response to lowercase for case-insensitive matching
	result = strings.ToLower(result)

	// Check for indicators that domain is definitely registered
	registeredIndicators := []string{
		"registrar:",
		"registrant:",
		"creation date:",
		"updated date:",
		"expiration date:",
		"name server:",
		"nserver:",
	}

	// Check for indicators that domain is definitely available
	availableIndicators := []string{
		"no match for",
		"not found",
		"no data found",
		"no entries found",
		"domain not found",
		"no object found",
		"no matching record",
		"status: free",
		"status: available",
	}

	// First check if domain is definitely available
	for _, indicator := range availableIndicators {
		if strings.Contains(result, indicator) {
			return true, nil
		}
	}

	// Then check if domain is definitely registered
	for _, indicator := range registeredIndicators {
		if strings.Contains(result, indicator) {
			return false, nil
		}
	}

	// If we can't determine the status from WHOIS and no DNS records exist,
	// assume the domain is available
	return true, nil
}

func generateDomains(length int, suffix string, pattern string, regexFilter string) []string {
	var domains []string
	letters := "abcdefghijklmnopqrstuvwxyz"
	numbers := "0123456789"

	// Compile regex if provided
	var regex *regexp.Regexp
	var err error
	if regexFilter != "" {
		regex, err = regexp.Compile(regexFilter)
		if err != nil {
			fmt.Printf("Invalid regex pattern: %v\n", err)
			os.Exit(1)
		}
	}

	switch pattern {
	case "d": // Pure numbers
		generateCombinations(&domains, "", numbers, length, suffix, regex)
	case "D": // Pure letters
		generateCombinations(&domains, "", letters, length, suffix, regex)
	case "a": // Alphanumeric
		generateCombinations(&domains, "", letters+numbers, length, suffix, regex)
	default:
		fmt.Println("Invalid pattern. Use -d for numbers, -D for letters, -a for alphanumeric")
		os.Exit(1)
	}

	return domains
}

func generateCombinations(domains *[]string, current string, charset string, length int, suffix string, regex *regexp.Regexp) {
	if len(current) == length {
		domain := current + suffix
		// Apply regex filter if provided
		if regex == nil || regex.MatchString(domain) {
			*domains = append(*domains, domain)
		}
		return
	}

	for _, c := range charset {
		generateCombinations(domains, current+string(c), charset, length, suffix, regex)
	}
}

func printHelp() {
	fmt.Println("Domain Scanner - A tool to check domain availability")
	fmt.Println("\nUsage:")
	fmt.Println("  go run main.go [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -l int      Domain length (default: 3)")
	fmt.Println("  -s string   Domain suffix (default: .li)")
	fmt.Println("  -p string   Domain pattern:")
	fmt.Println("              d: Pure numbers (e.g., 123.li)")
	fmt.Println("              D: Pure letters (e.g., abc.li)")
	fmt.Println("              a: Alphanumeric (e.g., a1b.li)")
	fmt.Println("  -r string   Regex filter for domain names")
	fmt.Println("  -delay int  Delay between queries in milliseconds (default: 1000)")
	fmt.Println("\nExamples:")
	fmt.Println("  1. Check 3-letter .li domains:")
	fmt.Println("     go run main.go -l 3 -s .li -p D")
	fmt.Println("\n  2. Check 3-digit .li domains:")
	fmt.Println("     go run main.go -l 3 -s .li -p d")
	fmt.Println("\n  3. Check 3-character alphanumeric .li domains:")
	fmt.Println("     go run main.go -l 3 -s .li -p a")
	fmt.Println("\n  4. Check domains containing 'abc':")
	fmt.Println("     go run main.go -l 5 -s .li -p D -r '.*abc.*'")
	fmt.Println("\n  5. Check domains starting with 'a' and ending with 'z':")
	fmt.Println("     go run main.go -l 4 -s .li -p D -r '^a.*z$'")
	fmt.Println("\n  6. Check domains containing only vowels:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -r '^[aeiou]+$'")
	fmt.Println("\n  7. Check domains with alternating letters and numbers:")
	fmt.Println("     go run main.go -l 4 -s .li -p a -r '^([a-z][0-9]){2}$'")
	fmt.Println("\nOutput files:")
	fmt.Println("  - available_domains_[pattern]_[length]_[suffix].txt")
	fmt.Println("  - registered_domains_[pattern]_[length]_[suffix].txt")
}

func main() {
	// Define command line flags
	length := flag.Int("l", 3, "Domain length")
	suffix := flag.String("s", ".li", "Domain suffix")
	pattern := flag.String("p", "D", "Domain pattern (d: numbers, D: letters, a: alphanumeric)")
	regexFilter := flag.String("r", "", "Regex filter for domain names")
	delay := flag.Int("delay", 1000, "Delay between queries in milliseconds")
	help := flag.Bool("h", false, "Show help information")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Ensure suffix starts with a dot
	if !strings.HasPrefix(*suffix, ".") {
		*suffix = "." + *suffix
	}

	domains := generateDomains(*length, *suffix, *pattern, *regexFilter)
	availableDomains := []string{}
	registeredDomains := []string{}

	fmt.Printf("Checking %d domains with pattern %s and length %d...\n", len(domains), *pattern, *length)
	if *regexFilter != "" {
		fmt.Printf("Using regex filter: %s\n", *regexFilter)
	}

	for i, domain := range domains {
		available, err := checkDomainAvailability(domain)
		if err != nil {
			fmt.Printf("[%d/%d] Error checking domain %s: %v\n", i+1, len(domains), domain, err)
			continue
		}

		if available {
			fmt.Printf("[%d/%d] Domain %s is AVAILABLE!\n", i+1, len(domains), domain)
			availableDomains = append(availableDomains, domain)
		} else {
			fmt.Printf("[%d/%d] Domain %s is REGISTERED\n", i+1, len(domains), domain)
			registeredDomains = append(registeredDomains, domain)
		}

		// Add delay to avoid rate limiting
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}

	// Save available domains to file
	availableFile := fmt.Sprintf("available_domains_%s_%d_%s.txt", *pattern, *length, strings.TrimPrefix(*suffix, "."))
	file, err := os.Create(availableFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, domain := range availableDomains {
		_, err := file.WriteString(domain + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
	}

	// Save registered domains to file
	registeredFile := fmt.Sprintf("registered_domains_%s_%d_%s.txt", *pattern, *length, strings.TrimPrefix(*suffix, "."))
	regFile, err := os.Create(registeredFile)
	if err != nil {
		fmt.Printf("Error creating registered domains file: %v\n", err)
		os.Exit(1)
	}
	defer regFile.Close()

	for _, domain := range registeredDomains {
		_, err := regFile.WriteString(domain + "\n")
		if err != nil {
			fmt.Printf("Error writing to registered domains file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nResults saved to:\n")
	fmt.Printf("- Available domains: %s\n", availableFile)
	fmt.Printf("- Registered domains: %s\n", registeredFile)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total domains checked: %d\n", len(domains))
	fmt.Printf("- Available domains: %d\n", len(availableDomains))
	fmt.Printf("- Registered domains: %d\n", len(registeredDomains))
}
