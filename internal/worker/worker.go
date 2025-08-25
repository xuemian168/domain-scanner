package worker

import (
	"time"

	"domain_scanner/internal/domain"
	"domain_scanner/internal/types"
)

func Worker(id int, jobs <-chan string, results chan<- types.DomainResult, delay time.Duration) {
	firstJob := true
	for domainName := range jobs {
		if firstJob {
			firstJob = false
			// This will be visible in domain check results
		}
		available, err := domain.CheckDomainAvailability(domainName)
		signatures, _ := domain.CheckDomainSignatures(domainName)
		results <- types.DomainResult{
			Domain:     domainName,
			Available:  available,
			Error:      err,
			Signatures: signatures,
		}
		time.Sleep(delay)
	}
}