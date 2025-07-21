package errors

import "fmt"

//go:generate stringer -type=EndpointsErrorCause
type EndpointsErrorCause int

const (
	_ EndpointsErrorCause = iota
	EndpointsRegionInvalid
	EndpointsRegionEmpty
	EndpointsStoredRegionInvalid
)

type EndpointsError struct {
	Region string
	Cause  EndpointsErrorCause
}

func (ee *EndpointsError) Error() string {
	return fmt.Sprintf("Cause: %s, Region: %s", ee.Cause.String(), ee.Region)
}
