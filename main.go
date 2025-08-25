package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"domain_scanner/internal/generator"
	"domain_scanner/internal/types"
	"domain_scanner/internal/worker"
)

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
	fmt.Println("  -regex-mode string Regex matching mode (default: full)")
	fmt.Println("    full: Match entire domain name")
	fmt.Println("    prefix: Match only domain name prefix")
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
	fmt.Println("\n  4. Use regex filter with full domain matching:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -r \"^[a-z]{2}[0-9]$\" -regex-mode full")
	fmt.Println("\n  5. Use regex filter with prefix matching:")
	fmt.Println("     go run main.go -l 3 -s .li -p D -r \"^[a-z]{2}\" -regex-mode prefix")
}

func calculateDomainsCount(pattern string, length int) int {
	var charsetSize int
	switch pattern {
	case "d":
		charsetSize = 10 // numbers 0-9
	case "D":
		charsetSize = 26 // letters a-z
	case "a":
		charsetSize = 36 // letters + numbers
	default:
		return 0
	}
	return int(math.Pow(float64(charsetSize), float64(length)))
}

func showMOTD() {
	fmt.Println("\033[1;36m") // Cyan color
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Domain Scanner v1.3.0                   ║")
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
	regexMode := flag.String("regex-mode", "full", "Regex match mode: 'full' or 'prefix'")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Ensure suffix starts with a dot
	if !strings.HasPrefix(*suffix, ".") {
		*suffix = "." + *suffix
	}

	// Determine regex mode
	var regexModeEnum types.RegexMode
	if *regexMode == "full" {
		regexModeEnum = types.RegexModeFull
	} else if *regexMode == "prefix" {
		regexModeEnum = types.RegexModePrefix
	} else {
		fmt.Println("Invalid regex-mode. Use 'full' or 'prefix'")
		os.Exit(1)
	}

	// Calculate total domains count and warn user if too many
	totalDomains := calculateDomainsCount(*pattern, *length)
	if totalDomains > 1000000 && *regexFilter == "" { // More than 1 million domains without regex filter
		fmt.Printf("\033[1;33mWarning: This configuration will generate %d domains.\033[0m\n", totalDomains)
		fmt.Printf("This may consume significant memory and time. Consider:\n")
		fmt.Printf("- Using shorter length (-l)\n")
		fmt.Printf("- Adding regex filter (-r)\n")
		fmt.Printf("- Using prefix matching (-regex-mode prefix)\n")
		fmt.Printf("\nContinue? (y/N): ")
		
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}
	} else if totalDomains > 1000000 && *regexFilter != "" {
		fmt.Printf("\033[1;33mInfo: Base configuration would generate %d domains, but regex filter will reduce this significantly.\033[0m\n", totalDomains)
	}

	domainsChan := generator.GenerateDomains(*length, *suffix, *pattern, *regexFilter, regexModeEnum)
	availableDomains := []string{}
	registeredDomains := []string{}

	// Display domain checking info
	if *regexFilter != "" {
		fmt.Printf("Checking domains with pattern %s and length %d using %d workers...\n", *pattern, *length, *workers)
		fmt.Printf("Using regex filter: %s (will show progress once generation completes)\n", *regexFilter)
	} else {
		fmt.Printf("Checking %d domains with pattern %s and length %d using %d workers...\n", totalDomains, *pattern, *length, *workers)
		fmt.Println("⚠️  Large domain set - progress tracking will be approximate")
	}

	// Create channels for jobs and results
	jobs := make(chan string, 1000)
	results := make(chan types.DomainResult, 1000)

	// Start workers
	for w := 1; w <= *workers; w++ {
		go worker.Worker(w, jobs, results, time.Duration(*delay)*time.Millisecond)
	}

	// Create a channel for domain status messages
	statusChan := make(chan string, 1000)

	// Start a goroutine to print status messages (above progress bar)
	go func() {
		for msg := range statusChan {
			fmt.Print("\r" + strings.Repeat(" ", 80) + "\r") // Clear progress bar
			fmt.Println(msg)                                   // Print message
		}
	}()

	// Track the number of active workers and jobs
	var jobsWG sync.WaitGroup
	domainsProcessed := 0

	// Create channels for coordination
	done := make(chan bool)
	totalGenerated := make(chan int, 1)
	
	// Start a goroutine to feed domains from generator to workers
	go func() {
		defer close(jobs)
		domainsGenerated := 0
		for domain := range domainsChan {
			domainsGenerated++
			jobsWG.Add(1)
			jobs <- domain
			if domainsGenerated%1000 == 0 {
				fmt.Printf("\rGenerated %d domains... (Workers processing in parallel)", domainsGenerated)
			}
		}
		fmt.Printf("\nTotal domains to check: %d\n", domainsGenerated)
		fmt.Println("Starting domain verification...")
		totalGenerated <- domainsGenerated
		
		// Wait for all jobs to complete, then signal completion
		jobsWG.Wait()
		close(results)
	}()

	// Collect results with progress tracking
	go func() {
		defer func() { done <- true }()
		defer close(statusChan)
		
		var total int
		totalKnown := false
		
		for result := range results {
			domainsProcessed++
			
			// Get total if available (non-blocking)
			if !totalKnown {
				select {
				case total = <-totalGenerated:
					totalKnown = true
				default:
					// Total not available yet
				}
			}
			
			// Create progress display
			var progress string
			var progressBar string
			if totalKnown && total > 0 {
				percentage := float64(domainsProcessed) / float64(total) * 100
				barWidth := 40
				filled := int(percentage / 100.0 * float64(barWidth))
				bar := ""
				for i := 0; i < barWidth; i++ {
					if i < filled {
						bar += "█"
					} else {
						bar += "░"
					}
				}
				progress = fmt.Sprintf("[%d/%d]", domainsProcessed, total)
				progressBar = fmt.Sprintf("\rProgress: %s [%s] %.1f%%", progress, bar, percentage)
			} else if totalDomains > 0 && *regexFilter == "" {
				// For large datasets without regex, show approximate progress
				percentage := float64(domainsProcessed) / float64(totalDomains) * 100
				barWidth := 40
				filled := int(percentage / 100.0 * float64(barWidth))
				bar := ""
				for i := 0; i < barWidth; i++ {
					if i < filled {
						bar += "█"
					} else {
						bar += "░"
					}
				}
				progress = fmt.Sprintf("[%d/~%d]", domainsProcessed, totalDomains)
				progressBar = fmt.Sprintf("\rProgress: %s [%s] ~%.1f%%", progress, bar, percentage)
			} else {
				progress = fmt.Sprintf("[%d]", domainsProcessed)
				progressBar = fmt.Sprintf("\rProgress: %s [Generating domains...]", progress)
			}
			
			if result.Error != nil {
				statusChan <- fmt.Sprintf("%s Error checking domain %s: %v", progress, result.Domain, result.Error)
			} else if result.Available {
				statusChan <- fmt.Sprintf("%s Domain %s is AVAILABLE!", progress, result.Domain)
				availableDomains = append(availableDomains, result.Domain)
			} else if *showRegistered {
				sigStr := strings.Join(result.Signatures, ", ")
				statusChan <- fmt.Sprintf("%s Domain %s is REGISTERED [%s]", progress, result.Domain, sigStr)
				registeredDomains = append(registeredDomains, result.Domain)
			}
			
			// Print progress bar
			fmt.Print(progressBar)
			
			jobsWG.Done()
		}
	}()
	
	// Wait for completion
	<-done
	
	// Clear progress bar and show completion
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
	fmt.Printf("✅ Domain checking completed!\n\n")

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
	fmt.Printf("- Total domains checked: %d\n", domainsProcessed)
	fmt.Printf("- Available domains: %d\n", len(availableDomains))
	if *showRegistered {
		fmt.Printf("- Registered domains: %d\n", len(registeredDomains))
	}
}
