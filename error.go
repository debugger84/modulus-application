package application

import "context"

type ActionErrorKey string

const WrongRequestDecoding ActionErrorKey = "WrongRequestDecoding"
const InvalidRequest ActionErrorKey = "InvalidRequest"

type ActionError struct {
	Ctx              context.Context
	Key              ActionErrorKey
	Err              error
	ValidationErrors []ValidationError
}

func (e *ActionError) Error() string {
	return e.Err.Error()
}

type ValidationError struct {
	field string
	err   string
}

func (e ValidationError) Error() string {
	return e.err
}

func NewServerError(ctx context.Context, key ActionErrorKey, err error) ActionResponse {
	return ActionResponse{
		StatusCode: 500,
		Error: &ActionError{
			Ctx:              ctx,
			Key:              key,
			Err:              err,
			ValidationErrors: nil,
		},
	}
}

func NewValidationError(ctx context.Context, errors []ValidationError) ActionResponse {
	if len(errors) == 0 {
		errors[0] = ValidationError{
			err: "Unknown error",
		}
	}
	return ActionResponse{
		StatusCode: 400,
		Error: &ActionError{
			Ctx:              ctx,
			Key:              InvalidRequest,
			Err:              errors[0],
			ValidationErrors: errors,
		},
	}
}
