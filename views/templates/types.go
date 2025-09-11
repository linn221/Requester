package templates

// FilterState holds the current state of filters applied to the requests list
type FilterState struct {
	Search      string
	ImportJobID string
	EndpointIDs []string
	Methods     []string
	Statuses    []string
	Types       []string
	SizeMin     string
	SizeMax     string
	Orders      []OrderClause
}

// OrderClause represents a single order clause
type OrderClause struct {
	Column    string
	Direction string
}

// ProgramOption represents a program option for dropdowns
type ProgramOption struct {
	ID   uint
	Name string
}

// EndpointOption represents an endpoint option for dropdowns
type EndpointOption struct {
	ID       uint
	Method   string
	Domain   string
	URI      string
	FullPath string
}

// ImportSummary represents the summary of an import operation
type ImportSummary struct {
	TotalRequests   int
	UniqueEndpoints int
	UniqueDomains   int
	Methods         map[string]int
	StatusCodes     map[string]int
}
