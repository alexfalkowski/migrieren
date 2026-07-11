package migrate

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type invalidConfigError struct {
	err error
}

func (e *invalidConfigError) Error() string {
	return ErrInvalidConfig.Error()
}

func (e *invalidConfigError) Unwrap() error {
	return ErrInvalidConfig
}

func (e *invalidConfigError) Stage() string {
	staged, ok := errors.AsType[stagedError](e.err)
	if !ok {
		return strings.Empty
	}

	return staged.Stage()
}

type stagedError interface {
	error
	Stage() string
}
