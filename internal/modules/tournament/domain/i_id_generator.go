package domain

// Interface used for avoding breaking a dependency of direction principle
// of modularity in terms of using a common utility reused among different modules.
type IDGenerator interface {
	GenerateID() string
	GenerateIDWithPrefixes(prefix1 string, prefix2 string) string
}
