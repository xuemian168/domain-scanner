package domain

import (
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"time"

	"domain_scanner/internal/reserved"

	"github.com/likexian/whois"
)

var (
	// Pre-initialized maps for O(1) lookup
	availableIndicatorsMap   map[string]bool
	unavailableIndicatorsMap map[string]bool
	indicatorsOnce           sync.Once

	// WHOIS servers for domain availability checking
	whoisServers = []string{
		"", // Default flow (IANA lookup)
		"whois.nic.li:43",
		"whois.nic.cx:43",        // Christmas Island
		"whois.nic.cz:43",        // Czech Republic
		"whois.verisign-grs.com:43",
		"whois.porkbun.com:43",
		"whois.godaddy.com:43",
		"whois.internic.net:43",
	}

	// WHOIS indicators for domain status detection
	registeredIndicators = []string{
		"registrar:",
		"registrant:",
		"creation date:",
		"updated date:",
		"expiration date:",
		"name server:",
		"nserver:",
		"status: connect",
		"changed:",
	}

	reservedIndicators = []string{
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
		"status: premium",
		"status: premium domain",
		"status: reserved by registry",
		"status: reserved by registrar",
		"status: reserved by administrator",
		"status: reserved by sponsoring organization",
		"status: reserved by iana",
		"status: reserved by icann",
		"status: trademark protected",
		"status: trademark reservation",
		"status: brand protection",
		"status: dpml block",
		"status: sunrise block",
		"status: landrush block",
		"status: hold",
		"status: frozen",
		"status: locked",
		"status: suspended",
		"status: quarantine",
		"status: redemption",
		"status: grace period",
		"status: pending delete",
		"status: pending restore",
		"status: clienthold",
		"status: serverhold",
		"status: clienttransferprohibited",
		"status: servertransferprohibited",
		"status: clientupdateprohibited",
		"status: serverupdateprohibited",
		"status: clientdeleteprohibited",
		"status: serverdeleteprohibited",
		"status: clientrenewprohibited",
		"status: serverrenewprohibited",
		"registry reserved",
		"registrar reserved",
		"reserved by",
		"reserved for",
		"reserved domain",
		"reserved name",
		"premium domain",
		"premium name",
		"trademark protected",
		"trademark block",
		"brand protection",
		"policy reserved",
		"policy block",
		"regulatory reserved",
		"regulatory block",
		"unavailable for registration",
		"not available for public registration",
		"not available for general registration",
		"registration not permitted",
		"registration prohibited",
		"registration restricted",
		"registration blocked",
		"registration suspended",
		"registration reserved",
		"this domain is reserved",
		"this name is reserved",
		"domain reserved",
		"name reserved",
		"domain blocked",
		"name blocked",
		"domain restricted",
		"name restricted",
		"domain unavailable",
		"name unavailable",
		"domain not available",
		"name not available",
		"domain withheld",
		"name withheld",
		"domain protected",
		"name protected",
		"domain frozen",
		"name frozen",
		"domain locked",
		"name locked",
		"domain suspended",
		"name suspended",
		"domain quarantined",
		"name quarantined",
		"domain on hold",
		"name on hold",
		"domain in grace period",
		"name in grace period",
		"domain pending delete",
		"name pending delete",
		"domain pending restore",
		"name pending restore",
	}

	// WHOIS indicators for domain availability detection
	availableIndicators = []string{
		"no match for", "not found", "no data found", "no entries found",
		"domain not found", "no object found", "no matching record",
		"status: free", "status: available", "is available for registration",
		"domain status: no object found", "no match!!", "not registered",
		"available for registration", "domain available", "available domain",
		"free domain", "domain free", "unregistered domain", "domain unregistered",
		"no match", "not found in database", "no matching record found",
		"domain name not found", "object does not exist", "no such domain",
		"domain status: available", "registration status: available",
		"state: available", "domain state: available", "available for purchase",
		"this domain is available", "domain is available", "can be registered",
		"eligible for registration", "free for registration", "open for registration",
		"ready for registration", "registration available", "status code: 210",
		"status code: 220", "response: 210", "response: 220",
		// .cx specific indicators
		"the queried object does not exist: no object found",
		"the queried object does not exist",
		// .cz specific indicators
		"%error:101: no entries found",
		"%error:101",
	}

	// Error indicators that should NOT be treated as "available"
	// These indicate service issues, not domain availability
	serviceErrorIndicators = []string{
		"the domain name search is temporarily unavailable",
		"temporarily unavailable",
		"service unavailable",
		"please try again later",
		"requests of this client are not permitted",
		"too many requests",
		"rate limit exceeded",
		"query limit exceeded",
		"access denied",
		"connection timeout",
		"service timeout",
	}

	unavailableIndicators = []string{
		"registrar:", "registrant:", "creation date:", "updated date:",
		"expiration date:", "name server:", "nserver:", "status: registered",
		"status: active", "status: ok", "status: connect",
		"status: clienttransferprohibited", "status: servertransferprohibited",
		"domain status: registered", "domain status: active", "registration date:",
		"expiry date:", "registry expiry date:", "registrar registration expiration date:",
		"admin contact:", "tech contact:", "billing contact:", "dnssec:",
		"domain servers in listed order:", "registered domain", "registered on:",
		"expires on:", "last updated on:", "changed:", "holder:", "person:",
		"sponsoring registrar:", "whois server:", "referral url:",
		"registry domain id:", "registrar whois server:", "registrar url:",
		"registrar iana id:", "registrar abuse contact email:",
		"registrar abuse contact phone:", "reseller:", "domain status:",
		"dnssec: unsigned", "dnssec: signed",
		// .cx specific indicators
		"registry expiry date:",
		"domain status: active",
		"registrar url:",
		// .cz specific indicators
		"registered:",
		"expire:",
		"nsset:",
		"admin-c:",
	}
)

// initIndicatorMaps initializes the indicator maps for fast lookup
func initIndicatorMaps() {
	indicatorsOnce.Do(func() {
		// Initialize available indicators map
		availableIndicatorsMap = make(map[string]bool, len(availableIndicators))
		for _, indicator := range availableIndicators {
			availableIndicatorsMap[indicator] = true
		}

		// Initialize unavailable indicators map
		unavailableIndicatorsMap = make(map[string]bool, len(unavailableIndicators))
		for _, indicator := range unavailableIndicators {
			unavailableIndicatorsMap[indicator] = true
		}
	})
}

func CheckDomainSignatures(domain string) ([]string, error) {
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
	maxRetries := 3
	baseDelay := 2 * time.Second

	for _, server := range whoisServers {
		for i := 0; i < maxRetries; i++ {
			var result string
			var err error
			if server == "" {
				result, err = whois.Whois(domain)
			} else {
				result, err = whois.Whois(domain, server)
			}

			if err == nil && result != "" {
				resultLower := strings.ToLower(result)
				// Check for registered indicators
				for _, indicator := range registeredIndicators {
					if strings.Contains(resultLower, indicator) {
						signatures = append(signatures, "WHOIS")
						return signatures, nil // Found registered, return immediately
					}
				}

				// Check for reserved indicators
				for _, indicator := range reservedIndicators {
					if strings.Contains(resultLower, indicator) {
						signatures = append(signatures, "RESERVED")
						return signatures, nil // Found reserved, return immediately
					}
				}

				// If we get here, the result was unclear, try next server
				break
			}

			// If there are still retry attempts, use exponential backoff
			if i < maxRetries-1 {
				// Calculate exponential delay: baseDelay * 2^i
				delay := baseDelay * time.Duration(1<<i)
				time.Sleep(delay)
			}
		}
		if server != "" {
			time.Sleep(1 * time.Second) // Delay before trying next server
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

func CheckDomainAvailability(domain string) (bool, error) {
	// First check if domain is reserved by pattern or TLD rules
	if reserved.IsReservedDomain(domain) {
		return false, nil
	}

	signatures, err := CheckDomainSignatures(domain)
	if err != nil {
		return false, err
	}

	// Check for reserved signature
	for _, sig := range signatures {
		if sig == "RESERVED" {
			return false, nil
		}
	}

	// If any other signatures found, domain is registered
	if len(signatures) > 0 {
		return false, nil
	}

	// Final WHOIS check for availability
	return checkWHOISAvailability(domain)
}

func checkWHOISAvailability(domain string) (bool, error) {
	maxRetries := 3
	baseDelay := 2 * time.Second
	foundAnyResult := false

	for _, server := range whoisServers {
		for i := 0; i < maxRetries; i++ {
			var result string
			var err error
			if server == "" {
				result, err = whois.Whois(domain)
			} else {
				result, err = whois.Whois(domain, server)
			}

			if err == nil && result != "" {
				foundAnyResult = true
				resultLower := strings.ToLower(result)

				// FIRST: Check for service errors (should NOT be treated as "available")
				if isServiceError(resultLower) {
					// Service error - treat as unavailable to prevent false positives
					return false, nil
				}

				// SECOND: Check for available indicators
				// Only return true if we have explicit "available" signal
				if isAvailableFromWHOIS(resultLower) {
					return true, nil
				}

				// THIRD: Check for unavailable indicators (check both original and lowercase)
				if isUnavailableFromWHOIS(result) || isUnavailableFromWHOIS(resultLower) {
					return false, nil
				}
				break // Move to next server if result is unclear
			}

			// If there are still retry attempts, use exponential backoff
			if i < maxRetries-1 {
				// Calculate exponential delay: baseDelay * 2^i
				delay := baseDelay * time.Duration(1<<i)
				time.Sleep(delay)
			}
		}
		if server != "" {
			time.Sleep(1 * time.Second) // Delay before trying next server
		}
	}

	// CRITICAL CHANGE: Conservative approach to prevent false positives
	// Default to UNAVAILABLE if no clear indication
	// Better to miss a potentially available domain than to report a registered one as available
	if !foundAnyResult {
		// No WHOIS data could be retrieved - assume domain is NOT available
		return false, nil
	}

	// WHOIS data was retrieved but couldn't determine status
	// Apply conservative approach: assume NOT available to prevent false positives
	return false, nil
}

func isAvailableFromWHOIS(result string) bool {
	// Most common patterns first for early return
	if strings.Contains(result, "status: free") ||
		strings.Contains(result, "not found") ||
		strings.Contains(result, "no match") ||
		strings.Contains(result, "status: available") ||
		strings.Contains(result, "no data found") ||
		strings.Contains(result, "is available") {
		return true
	}

	// Less common patterns
	initIndicatorMaps()
	for indicator := range availableIndicatorsMap {
		if strings.Contains(result, indicator) {
			return true
		}
	}

	return false
}

func isUnavailableFromWHOIS(result string) bool {
	// Most common patterns first for early return
	if strings.Contains(result, "registrar:") ||
		strings.Contains(result, "name server:") ||
		strings.Contains(result, "nserver:") ||
		strings.Contains(result, "creation date:") ||
		strings.Contains(result, "status: connect") ||
		strings.Contains(result, "Nserver:") ||
		strings.Contains(result, "Changed:") {
		return true
	}

	// Less common patterns
	initIndicatorMaps()
	for indicator := range unavailableIndicatorsMap {
		if strings.Contains(result, indicator) {
			return true
		}
	}

	return false
}

func isServiceError(result string) bool {
	// Check for service error indicators that should NOT be treated as "available"
	for _, indicator := range serviceErrorIndicators {
		if strings.Contains(result, indicator) {
			return true
		}
	}
	return false
}
