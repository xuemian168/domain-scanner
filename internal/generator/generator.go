package generator

import (
	"fmt"
	"os"
	"regexp"

	"domain_scanner/internal/types"
)

func GenerateDomains(length int, suffix string, pattern string, regexFilter string, regexMode types.RegexMode) []string {
	var domains []string
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

	switch pattern {
	case "d":
		generateCombinations(&domains, "", numbers, length, suffix, regex, regexMode)
	case "D":
		generateCombinations(&domains, "", letters, length, suffix, regex, regexMode)
	case "a":
		generateCombinations(&domains, "", letters+numbers, length, suffix, regex, regexMode)
	default:
		fmt.Println("Invalid pattern. Use -d for numbers, -D for letters, -a for alphanumeric")
		os.Exit(1)
	}

	return domains
}

func generateCombinations(domains *[]string, current string, charset string, length int, suffix string, regex *regexp.Regexp, regexMode types.RegexMode) {
	if len(current) == length {
		domain := current + suffix
		var match bool
		switch regexMode {
		case types.RegexModeFull:
			match = regex == nil || regex.MatchString(domain)
		case types.RegexModePrefix:
			match = regex == nil || regex.MatchString(current)
		}

		if match {
			*domains = append(*domains, domain)
		}
		return
	}

	for _, c := range charset {
		generateCombinations(domains, current+string(c), charset, length, suffix, regex, regexMode)
	}
}