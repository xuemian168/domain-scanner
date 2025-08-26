package reserved

import (
	"regexp"
	"strings"
	"sync"
)

var (
	// Pre-compiled regex patterns for better performance
	reservedPatternsOnce sync.Once
	compiledPatterns     []*regexp.Regexp
	reservedWordsMap     map[string]bool
	techPrefixesMap      map[string]bool
)

// initReservedPatterns initializes the compiled patterns and word maps
func initReservedPatterns() {
	reservedPatternsOnce.Do(func() {
		// Simple patterns that can be compiled
		patterns := []string{
			"^[a-z]$",      // Single letter
			"^[a-z]{2}$",   // Two letters
			"^[0-9]{2,3}$", // 2-3 digits
			"^.{1,2}$",     // Very short domains
		}

		for _, pattern := range patterns {
			if re, err := regexp.Compile(pattern); err == nil {
				compiledPatterns = append(compiledPatterns, re)
			}
		}

		// Initialize reserved words map for O(1) lookup
		initReservedWordsMap()

		// Initialize tech prefixes map
		initTechPrefixesMap()
	})
}

// initReservedWordsMap creates a map of reserved words for fast lookup
func initReservedWordsMap() {
	// Common reserved words
	commonWords := []string{
		"www", "ftp", "mail", "email", "smtp", "pop", "imap", "ns", "dns", "mx",
		"admin", "root", "test", "demo", "example", "localhost", "api", "app", "web",
		"site", "blog", "shop", "store", "com", "net", "org", "gov", "edu", "mil", "int",
		"info", "biz", "name", "pro", "museum", "coop", "aero", "jobs", "mobi", "travel",
		"xxx", "tel", "asia", "cat", "post", "geo",
	}

	// Common service names (sample set - can be extended)
	serviceNames := []string{
		"google", "facebook", "twitter", "youtube", "amazon", "microsoft", "apple",
		"netflix", "instagram", "linkedin", "whatsapp", "telegram", "github", "gitlab",
		"bitbucket", "stackoverflow", "reddit", "wikipedia", "cloudflare", "aws", "azure",
		"docker", "kubernetes", "nginx", "apache", "mysql", "postgresql", "mongodb", "redis",
		"stripe", "paypal", "bitcoin", "ethereum", "wordpress", "shopify", "zoom", "slack",
		// Add more as needed...
	}

	// Common generic terms
	genericTerms := []string{
		"login", "register", "signup", "signin", "logout", "profile", "account", "dashboard",
		"settings", "config", "preferences", "privacy", "security", "terms", "conditions",
		"policy", "legal", "help", "support", "contact", "about", "faq", "blog", "news",
		"press", "media", "careers", "jobs", "team", "company", "home", "index", "main",
		"default", "landing", "welcome", "hello", "start", "begin", "download", "upload",
		"search", "find", "discover", "explore", "browse", "navigate", "menu", "navbar",
		// Add more as needed...
	}

	// Combine all words
	allWords := append(append(commonWords, serviceNames...), genericTerms...)

	reservedWordsMap = make(map[string]bool, len(allWords))
	for _, word := range allWords {
		reservedWordsMap[word] = true
	}
}

// initTechPrefixesMap creates a map of technical prefixes for fast lookup
func initTechPrefixesMap() {
	techPrefixes := []string{
		"localhost", "dns", "ns", "mx", "mail", "smtp", "pop", "imap", "ftp",
		"www", "web", "server", "host", "node", "db", "cache", "cdn", "api",
		"app", "admin", "root", "sys", "net", "org", "gov", "edu", "mil",
		"int", "com", "info", "biz", "name", "pro",
	}

	techPrefixesMap = make(map[string]bool, len(techPrefixes))
	for _, prefix := range techPrefixes {
		techPrefixesMap[prefix] = true
	}
}

// IsReservedByPattern checks if a domain is reserved based on common patterns
func IsReservedByPattern(domain string) bool {
	initReservedPatterns()

	domainLower := strings.ToLower(domain)

	// Remove TLD to check only the domain name part
	parts := strings.Split(domainLower, ".")
	if len(parts) < 2 {
		return false
	}
	domainName := parts[0]

	// Check with map first (O(1) lookup)
	if reservedWordsMap[domainName] {
		return true
	}

	// Check compiled regex patterns
	for _, re := range compiledPatterns {
		if re.MatchString(domainName) {
			return true
		}
	}

	// Check technical terms with numbers (special pattern)
	if checkTechnicalPattern(domainName) {
		return true
	}

	// Check IP-like patterns
	if domainName == "127" || domainName == "192" || domainName == "10" ||
		domainName == "172" || domainName == "255" {
		return true
	}

	return false
}

// checkTechnicalPattern checks for technical terms with optional numbers
func checkTechnicalPattern(name string) bool {
	// Direct match first
	if techPrefixesMap[name] {
		return true
	}

	// Check for patterns with number suffix
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] < '0' || name[i] > '9' {
			// Found non-digit, check if prefix is technical term
			if i < len(name)-1 {
				prefix := name[:i+1]
				if techPrefixesMap[prefix] {
					return true
				}
			}
			break
		}
	}

	return false
}

// Cache for TLD rules to avoid repeated map creation
var (
	tldRulesCache     map[string]map[string]bool
	tldRulesCacheOnce sync.Once
)

// initTLDRulesCache initializes the TLD rules cache
func initTLDRulesCache() {
	tldRulesCacheOnce.Do(func() {
		tldRulesCache = make(map[string]map[string]bool)

		// Define TLD-specific reserved domains
		tldRules := map[string][]string{
			".com": {
				"com", "net", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "google", "facebook", "twitter", "youtube", "amazon", "microsoft",
				"apple", "netflix", "instagram", "linkedin",
			},
			".net": {
				"net", "com", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "network", "internet", "intranet", "extranet", "lan", "wan", "vpn",
			},
			".org": {
				"org", "com", "net", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "organization", "foundation", "charity", "nonprofit", "ngo",
			},
			".li": {
				"li", "com", "net", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "liechtenstein", "principality", "government", "official", "royal",
			},
			".io": {
				"io", "com", "net", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "input", "output", "tech", "technology", "startup", "developer",
			},
			".ai": {
				"ai", "com", "net", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "artificial", "intelligence", "machine", "learning", "neural", "deep",
			},
			".de": {
				"de", "com", "net", "org", "edu", "gov", "mil", "int", "www", "ftp", "mail",
				"email", "smtp", "pop", "imap", "dns", "ns", "mx", "web", "site", "blog",
				"shop", "store", "app", "api", "admin", "root", "test", "demo", "example",
				"localhost", "deutschland", "german", "germany", "berlin", "munich", "hamburg",
			},
			// Add more TLDs as needed...
		}

		// Convert to maps for O(1) lookup
		for tld, domains := range tldRules {
			domainMap := make(map[string]bool, len(domains))
			for _, domain := range domains {
				domainMap[domain] = true
			}
			tldRulesCache[tld] = domainMap
		}
	})
}

// IsReservedByTLD checks if a domain is reserved based on TLD-specific rules
func IsReservedByTLD(domain string) bool {
	initTLDRulesCache()

	domainLower := strings.ToLower(domain)

	// Extract TLD from domain
	parts := strings.Split(domainLower, ".")
	if len(parts) < 2 {
		return false
	}

	tld := "." + parts[len(parts)-1]
	domainName := parts[0]

	// Check if domain name is in the reserved list for this TLD
	if reservedMap, exists := tldRulesCache[tld]; exists {
		return reservedMap[domainName]
	}

	return false
}

// IsReservedDomain checks if a domain is reserved using multiple methods
func IsReservedDomain(domain string) bool {
	// Check pattern-based rules
	if IsReservedByPattern(domain) {
		return true
	}

	// Check TLD-specific rules
	if IsReservedByTLD(domain) {
		return true
	}

	return false
}
