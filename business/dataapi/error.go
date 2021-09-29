package dataapi

// Error implements the error interface and is used to identify a trusted error
// The DataAPI uses a cascading error mechanism to report on multiple errors
// Thus this struct is not used and is here just to complete the architecture
type Error struct {
	Err error
}

func (err *Error) Error() string {
	return err.Err.Error()
}
