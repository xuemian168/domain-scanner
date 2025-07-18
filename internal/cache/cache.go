package cache

import (
	"sync"
	"time"
)

// DomainCache is a thread-safe cache for domain check results
type DomainCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheEntry
	ttl   time.Duration
}

// CacheEntry represents a cached domain check result
type CacheEntry struct {
	Available  bool
	Signatures []string
	Timestamp  time.Time
}

// NewDomainCache creates a new domain cache with specified TTL
func NewDomainCache(ttl time.Duration) *DomainCache {
	return &DomainCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}
}

// Get retrieves a cached result if it exists and is not expired
func (dc *DomainCache) Get(domain string) (available bool, signatures []string, found bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	
	entry, exists := dc.cache[domain]
	if !exists {
		return false, nil, false
	}
	
	// Check if entry is expired
	if time.Since(entry.Timestamp) > dc.ttl {
		return false, nil, false
	}
	
	return entry.Available, entry.Signatures, true
}

// Set stores a domain check result in the cache
func (dc *DomainCache) Set(domain string, available bool, signatures []string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	
	dc.cache[domain] = &CacheEntry{
		Available:  available,
		Signatures: signatures,
		Timestamp:  time.Now(),
	}
}

// Clean removes expired entries from the cache
func (dc *DomainCache) Clean() {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	
	now := time.Now()
	for domain, entry := range dc.cache {
		if now.Sub(entry.Timestamp) > dc.ttl {
			delete(dc.cache, domain)
		}
	}
}

// StartCleanupRoutine starts a background routine to clean expired entries
func (dc *DomainCache) StartCleanupRoutine(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			dc.Clean()
		}
	}()
}