package render

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Possible bucket class names
	Classes = []string{"standard", "vault", "cold", "flex", "smart", "onerate_active"}
	// Class tokenizer
	classTokenizer = strings.Join(Classes, `|`)
	/// Region specifications
	regionSpec = fmt.Sprintf(`(?i)^(\w+(?:-\w+)??)(-geo)?(?:-(%s|\*))?$`, classTokenizer)
	// RegionDecoderRegex regular expression to break down regions and storage class
	RegionDecoderRegex = regexp.MustCompile(regionSpec)
)

// Retrieve class from the location constraint
func getClassFromLocationConstraint(location string) string {
	// Break down regions and storage class
	regionDetails := RegionDecoderRegex.FindStringSubmatch(location)
	if regionDetails != nil {
		return regionDetails[3]
	}
	// Return no class
	return ""
}

// Retrieve region from the location constraint
func getRegionFromLocationConstraint(location string) string {
	// Break down regions and storage class
	regionDetails := RegionDecoderRegex.FindStringSubmatch(location)
	if regionDetails != nil {
		return regionDetails[1]
	}
	// Return no region
	return ""
}

// Render class helper
func renderClass(class string) string {
	switch class {
	case "":
		return "Standard"
	case "cold":
		return "Cold Vault"
	case "onerate_active":
		return "One-rate Active"
	default:
		return strings.Title(class)
	}
}

//
// FormatFileSize function is excerpted from from
// the following link:
// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
// Description: it outputs a human readable representation of the value using
// multiples of 1024
func FormatFileSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
