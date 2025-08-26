package main

import (
	"flag"
	"fmt"
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

func showMOTD() {
	fmt.Println("\033[1;36m") // Cyan color
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Domain Scanner v1.3.2                   ║")
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

	domainChan := generator.GenerateDomains(*length, *suffix, *pattern, *regexFilter, regexModeEnum)
	availableDomains := []string{}
	registeredDomains := []string{}

	// 计算总域名数量
	totalDomains := calculateDomainsCount(*length, *pattern)
	fmt.Printf("Checking %d domains with pattern %s and length %d using %d workers...\n",
		totalDomains, *pattern, *length, *workers)
	if *regexFilter != "" {
		fmt.Printf("Using regex filter: %s\n", *regexFilter)
	}

	// Create channels for jobs and results
	jobs := make(chan string, 1000)
	results := make(chan types.DomainResult, 1000)

	// Start workers
	for w := 1; w <= *workers; w++ {
		go worker.Worker(w, jobs, results, time.Duration(*delay)*time.Millisecond)
	}

	// Send jobs from domain generator
	go func() {
		defer close(jobs)
		for domain := range domainChan {
			jobs <- domain
		}
	}()

	// Create a channel for domain status messages
	statusChan := make(chan string, 1000)

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
		processedCount := 0
		for result := range results {
			processedCount++
			progress := fmt.Sprintf("[%d/%d]", processedCount, totalDomains)
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

			// 检查是否已处理完所有域名
			if processedCount >= totalDomains {
				break
			}
		}
		close(statusChan)
	}()

	// 监控任务完成
	go func() {
		// 等待所有域名生成完毕
		for range domainChan {
			// domainChan 关闭后此循环结束
		}
		// 等待一段时间确保所有结果都被处理
		time.Sleep(1 * time.Second)
		close(results)
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
	fmt.Printf("- Total domains checked: %d\n", totalDomains)
	fmt.Printf("- Available domains: %d\n", len(availableDomains))
	if *showRegistered {
		fmt.Printf("- Registered domains: %d\n", len(registeredDomains))
	}
}

// calculateDomainsCount 计算给定模式和长度的域名总数
func calculateDomainsCount(length int, pattern string) int {
	var charsetSize int
	switch pattern {
	case "d": // 纯数字
		charsetSize = 10 // 0-9
	case "D": // 纯字母
		charsetSize = 26 // a-z
	case "a": // 字母数字
		charsetSize = 36 // a-z + 0-9
	default:
		return 0
	}

	total := 1
	for i := 0; i < length; i++ {
		total *= charsetSize
	}
	return total
}
