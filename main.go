package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
	"github.com/schollz/progressbar/v3"
)

type DomainResult struct {
	Domain    string
	Available bool
	Error     error
}

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

func worker(id int, jobs <-chan string, results chan<- DomainResult, delay time.Duration) {
	for domain := range jobs {
		available, err := checkDomainAvailability(domain)
		results <- DomainResult{
			Domain:    domain,
			Available: available,
			Error:     err,
		}
		time.Sleep(delay) // Add delay to avoid rate limiting
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

func main() {
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

	// Create progress bar
	bar := progressbar.NewOptions(len(domains),
		progressbar.OptionSetDescription("Scanning domains"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

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

	// Collect results
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(domains); i++ {
			result := <-results
			if result.Error != nil {
				fmt.Printf("\nError checking domain %s: %v\n", result.Domain, result.Error)
				bar.Add(1)
				continue
			}

			if result.Available {
				fmt.Printf("\nDomain %s is AVAILABLE!\n", result.Domain)
				availableDomains = append(availableDomains, result.Domain)
			} else if *showRegistered {
				fmt.Printf("\nDomain %s is REGISTERED\n", result.Domain)
				registeredDomains = append(registeredDomains, result.Domain)
			}
			bar.Add(1)
		}
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
