package generator

import (
	"fmt"
	"os"
	"regexp"

	"domain_scanner/internal/types"
)

// GenerateDomains returns a channel that streams domains instead of generating all at once
func GenerateDomains(length int, suffix string, pattern string, regexFilter string, regexMode types.RegexMode) <-chan string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	numbers := "0123456789"

	var regex *regexp.Regexp
	var err error
	if regexFilter != "" {
		regex, err = regexp.Compile(regexFilter)
		if err != nil {
			fmt.Printf("Invalid regex pattern: %v\n", err)
			os.Exit(1)
		}
	}

	domainChan := make(chan string, 1000) // Buffer for better performance

	go func() {
		defer close(domainChan)
		
		switch pattern {
		case "d":
			generateCombinationsIterative(domainChan, numbers, length, suffix, regex, regexMode)
		case "D":
			generateCombinationsIterative(domainChan, letters, length, suffix, regex, regexMode)
		case "a":
			generateCombinationsIterative(domainChan, letters+numbers, length, suffix, regex, regexMode)
		default:
			fmt.Println("Invalid pattern. Use -d for numbers, -D for letters, -a for alphanumeric")
			os.Exit(1)
		}
	}()

	return domainChan
}

// generateCombinationsIterative uses iterative approach instead of recursive to avoid stack overflow
func generateCombinationsIterative(domainChan chan<- string, charset string, length int, suffix string, regex *regexp.Regexp, regexMode types.RegexMode) {
	charsetSize := len(charset)
	if charsetSize == 0 || length <= 0 {
		return
	}

	// Use counter-based approach for generating combinations
	total := 1
	for i := 0; i < length; i++ {
		total *= charsetSize
	}

	for counter := 0; counter < total; counter++ {
		current := ""
		temp := counter
		
		// Generate domain string from counter
		for i := 0; i < length; i++ {
			current = string(charset[temp%charsetSize]) + current
			temp /= charsetSize
		}

		domain := current + suffix
		var match bool
		switch regexMode {
		case types.RegexModeFull:
			match = regex == nil || regex.MatchString(domain)
		case types.RegexModePrefix:
			match = regex == nil || regex.MatchString(current)
		}

		if match {
			domainChan <- domain
		}
	}
}