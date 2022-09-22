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

	// Check if the class is onerate_active, return "One-rate Active"
	if c == "onerate_active" {
		return "One-rate Active"
	}

	// Return
	return strings.Title(string(c))
}
