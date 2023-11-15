package service

import (
	"errors"
	"fmt"
)

type SocialNetworkAccountAlreadyExistsError struct {
	Message string
}

type PageAlreadyExistsError struct {
	Message string
}

type ValidationError struct {
	Message string
	Field   string
	Rule    string
}

type InternalError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

func NewSocialNetworkAccountAlreadyExistsError(message string) *SocialNetworkAccountAlreadyExistsError {
	return &SocialNetworkAccountAlreadyExistsError{
		Message: message,
	}
}

func NewPageAlreadyExistsError(message string) *PageAlreadyExistsError {
	return &PageAlreadyExistsError{
		Message: message,
	}
}

func NewValidationError(message string, field string, rule string) *ValidationError {
	return &ValidationError{
		Message: message,
		Field:   field,
		Rule:    rule,
	}
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		Message: message,
	}
}

func NewNotFoundError(message string) error {
	return &NotFoundError{
		Message: message,
	}
}

func IsSocialNetworkAccountAlreadyExistsError(err error) bool {
	var e *SocialNetworkAccountAlreadyExistsError

	return errors.As(err, &e)
}

func IsPageAlreadyExistsError(err error) bool {
	var e *PageAlreadyExistsError

	return errors.As(err, &e)
}

func IsValidationError(err error) bool {
	var e *ValidationError

	return errors.As(err, &e)
}

func IsInternalError(err error) bool {
	var e *InternalError

	return errors.As(err, &e)
}

func IsNotFoundError(err error) bool {
	var e *NotFoundError

	return errors.As(err, &e)
}

func (e *SocialNetworkAccountAlreadyExistsError) Error() string {
	return e.Message
}

func (e *PageAlreadyExistsError) Error() string {
	return e.Message
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"Error: %v, Field: %v, Rule: %v",
		e.Message,
		e.Field,
		e.Rule,
	)
}

func (e *InternalError) Error() string {
	return e.Message
}

func (e *NotFoundError) Error() string {
	return e.Message
}
