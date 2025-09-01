package types

type DomainResult struct {
	Domain     string
	Available  bool
	Error      error
	Signatures []string
}
