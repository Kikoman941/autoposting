package resolver

import (
	"github.com/ztrue/tracerr"
)

func NewResolverError(
	responseError string,
	originalError error,
) *ResolverError {
	return &ResolverError{
		tracerr.Errorf(responseError),
		originalError,
	}
}

type ResolverError struct { // nolint:revive
	ResponseError error
	OriginalError error
}

func (e *ResolverError) Error() string {
	return e.ResponseError.Error()
}
