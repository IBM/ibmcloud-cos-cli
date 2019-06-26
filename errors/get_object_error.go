//go:generate stringer -type=ObjectGetErrorCause

package errors

import "fmt"

type ObjectGetErrorCause int

const (
	_ ObjectGetErrorCause = iota
	IsDir
	Opening
)

type ObjectGetError struct {
	_                struct{}
	ParentError      error
	Location         string
	UsingDefaultRule bool
	Cause            ObjectGetErrorCause
}

func (oge *ObjectGetError) Error() string {
	return fmt.Sprintf("Cause: %s, Location %s, DefaultRule: %t, Parent: %s",
		oge.Cause.String(), oge.Location, oge.UsingDefaultRule, oge.ParentError)
}
