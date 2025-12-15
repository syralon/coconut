package errors

import (
	"fmt"
	"runtime/debug"
)

type RecoveryError struct {
	stack    string
	recovery error
}

func (r *RecoveryError) Error() string {
	return r.recovery.Error()
}

func (r *RecoveryError) Stack() string {
	return r.stack
}

func Recovery(rec any) *RecoveryError {
	if rec == nil {
		return nil
	}
	return &RecoveryError{
		stack:    string(debug.Stack()),
		recovery: any2error(rec),
	}
}

func any2error(v any) error {
	if err, ok := v.(error); ok {
		return err
	}
	return fmt.Errorf("%v", v)
}
