package functions

import "strings"

// DefaultClass variable = Standard
var DefaultClass = Class("standard")

// Class type as string
type Class string

// Helper function for class of a bucket
func (c Class) String() string {
	if c == "" {
		return strings.Title(string(DefaultClass))
	}

	// Check if the class is cold, append "Cold Vault"
	if c == "cold" {
		return "Cold Vault"
	}

	// Return
	return strings.Title(string(c))
}
