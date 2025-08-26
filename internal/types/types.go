package types

type DomainResult struct {
	Domain     string
	Available  bool
	Error      error
	Signatures []string
}

type RegexMode int

const (
	RegexModeFull RegexMode = iota
	RegexModePrefix
)
