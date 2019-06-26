package errors

// CodeError Interface allows classify and handle errors using their code
type CodeError interface {
	// Satisfy the generic error interface.
	error

	// Returns the Error Code
	Code() string
}
