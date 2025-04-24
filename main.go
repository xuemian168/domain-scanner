package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
)

type DomainResult struct {
	Domain     string
	Available  bool
	Error      error
	Signatures []string
}

func checkDomainSignatures(domain string) ([]string, error) {
	var signatures []string

	// 1. Check DNS NS records
	nsRecords, err := net.LookupNS(domain)
	if err == nil && len(nsRecords) > 0 {
		signatures = append(signatures, "DNS_NS")
	}

	// 2. Check DNS A records
	ipRecords, err := net.LookupIP(domain)
	if err == nil && len(ipRecords) > 0 {
		signatures = append(signatures, "DNS_A")
	}

	// 3. Check DNS MX records
	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		signatures = append(signatures, "DNS_MX")
	}

	// 4. Check WHOIS information with retry
	var whoisResult string
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		result, err := whois.Whois(domain)
		if err == nil {
			whoisResult = result
			break
		}
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2) // Wait 2 seconds before retry
		}
	}

	if whoisResult != "" {
		// Convert WHOIS response to lowercase for case-insensitive matching
		result := strings.ToLower(whoisResult)

		// Check for registration indicators
		registeredIndicators := []string{
			"registrar:",
			"registrant:",
			"creation date:",
			"updated date:",
			"expiration date:",
			"name server:",
			"nserver:",
		}

		for _, indicator := range registeredIndicators {
			if strings.Contains(result, indicator) {
				signatures = append(signatures, "WHOIS")
				break
			}
		}

		// Check for reserved domain indicators
		reservedIndicators := []string{
			"status: reserved",
			"status: restricted",
			"status: blocked",
			"status: prohibited",
			"status: reserved for registry",
			"status: reserved for registrar",
			"status: reserved for registry operator",
			"status: reserved for future use",
			"status: not available for registration",
			"status: not available for general registration",
			"status: reserved for special purposes",
			"status: reserved for government use",
			"status: reserved for educational institutions",
			"status: reserved for non-profit organizations",
		}

		for _, indicator := range reservedIndicators {
			if strings.Contains(result, indicator) {
				signatures = append(signatures, "RESERVED")
				break
			}
		}
	}

	// 5. Check SSL certificate with timeout
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err == nil {
		defer conn.Close()
		state := conn.ConnectionState()
		if len(state.PeerCertificates) > 0 {
			signatures = append(signatures, "SSL")
		}
	}

	return signatures, nil
}

func checkDomainAvailability(domain string) (bool, error) {
	signatures, err := checkDomainSignatures(domain)
	if err != nil {
		return false, err
	}

	// If domain is reserved, it's not available
	for _, sig := range signatures {
		if sig == "RESERVED" {
			return false, nil
		}
	}

	// If any other signature is found, domain is registered
	if len(signatures) > 0 {
		return false, nil
	}

	// If no signatures found, check WHOIS as final verification with retry
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		result, err := whois.Whois(domain)
		if err == nil {
			// Convert WHOIS response to lowercase for case-insensitive matching
			result = strings.ToLower(result)

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

			for _, indicator := range availableIndicators {
				if strings.Contains(result, indicator) {
					return true, nil
				}
			}
			break
		}
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2) // Wait 2 seconds before retry
		}
	}

	// If we can't determine the status, assume the domain is available
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

func worker(id int, jobs <-chan string, results chan<- DomainResult, delay time.Duration) {
	for domain := range jobs {
		available, err := checkDomainAvailability(domain)
		signatures, _ := checkDomainSignatures(domain)
		results <- DomainResult{
			Domain:     domain,
			Available:  available,
			Error:      err,
			Signatures: signatures,
		}
		time.Sleep(delay)
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
	fmt.Println("  -workers int Number of concurrent workers (default: 10)")
	fmt.Println("  -show-registered Show registered domains in output (default: false)")
	fmt.Println("  -h          Show help information")
	fmt.Println("\nExamples:")
	fmt.Println("  1. Check 3-letter .li domains with 20 workers:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -workers 20")
	fmt.Println("\n  2. Check domains with custom delay and workers:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -delay 500 -workers 15")
	fmt.Println("\n  3. Show both available and registered domains:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -show-registered")
}

func showMOTD() {
	fmt.Println("\033[1;36m") // Cyan color
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Domain Scanner v1.0                     ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  A powerful tool for checking domain name availability     ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  Developer: www.ict.run                                    ║")
	fmt.Println("║  GitHub:    https://github.com/xuemian168/domain-scanner   ║")
	fmt.Println("║                                                            ║")
	fmt.Println("║  License:   AGPL-3.0                                       ║")
	fmt.Println("║  Copyright © 2025                                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println("\033[0m") // Reset color
	fmt.Println()
}

func main() {
	// Show MOTD
	showMOTD()

	// Define command line flags
	length := flag.Int("l", 3, "Domain length")
	suffix := flag.String("s", ".li", "Domain suffix")
	pattern := flag.String("p", "D", "Domain pattern (d: numbers, D: letters, a: alphanumeric)")
	regexFilter := flag.String("r", "", "Regex filter for domain names")
	delay := flag.Int("delay", 1000, "Delay between queries in milliseconds")
	workers := flag.Int("workers", 10, "Number of concurrent workers")
	showRegistered := flag.Bool("show-registered", false, "Show registered domains in output")
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

	fmt.Printf("Checking %d domains with pattern %s and length %d using %d workers...\n",
		len(domains), *pattern, *length, *workers)
	if *regexFilter != "" {
		fmt.Printf("Using regex filter: %s\n", *regexFilter)
	}

	// Create channels for jobs and results
	jobs := make(chan string, len(domains))
	results := make(chan DomainResult, len(domains))

	// Start workers
	for w := 1; w <= *workers; w++ {
		go worker(w, jobs, results, time.Duration(*delay)*time.Millisecond)
	}

	// Send jobs
	for _, domain := range domains {
		jobs <- domain
	}
	close(jobs)

	// Create a channel for domain status messages
	statusChan := make(chan string, len(domains))

	// Start a goroutine to print status messages
	go func() {
		for msg := range statusChan {
			fmt.Println(msg)
		}
	}()

	// Collect results
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(domains); i++ {
			result := <-results
			progress := fmt.Sprintf("[%d/%d]", i+1, len(domains))
			if result.Error != nil {
				statusChan <- fmt.Sprintf("%s Error checking domain %s: %v", progress, result.Domain, result.Error)
				continue
			}

			if result.Available {
				statusChan <- fmt.Sprintf("%s Domain %s is AVAILABLE!", progress, result.Domain)
				availableDomains = append(availableDomains, result.Domain)
			} else if *showRegistered {
				sigStr := strings.Join(result.Signatures, ", ")
				statusChan <- fmt.Sprintf("%s Domain %s is REGISTERED [%s]", progress, result.Domain, sigStr)
				registeredDomains = append(registeredDomains, result.Domain)
			}
		}
		close(statusChan)
	}()
	wg.Wait()

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

	// Save registered domains to file only if show-registered is true
	registeredFile := fmt.Sprintf("registered_domains_%s_%d_%s.txt", *pattern, *length, strings.TrimPrefix(*suffix, "."))
	if *showRegistered {
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
	}

	fmt.Printf("\n\nResults saved to:\n")
	fmt.Printf("- Available domains: %s\n", availableFile)
	if *showRegistered {
		fmt.Printf("- Registered domains: %s\n", registeredFile)
	}
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total domains checked: %d\n", len(domains))
	fmt.Printf("- Available domains: %d\n", len(availableDomains))
	if *showRegistered {
		fmt.Printf("- Registered domains: %d\n", len(registeredDomains))
	}
}
