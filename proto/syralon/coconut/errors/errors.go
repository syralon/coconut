package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CodeUnknown = 0

	maxGRPCCode = 17
)

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *Error) Attrs() []slog.Attr {
	return []slog.Attr{
		slog.Int("code", int(e.GetCode())),
		slog.String("message", e.GetMessage()),
		slog.String("reason", e.GetReason()),
		slog.Int("status", int(e.GetStatusCode())),
	}
}

func (e *Error) Status() *status.Status {
	st := status.New(codes.Code(e.Code), e.Message)
	detail := &errdetails.ErrorInfo{
		Reason:   e.GetReason(),
		Metadata: map[string]string{"status_code": strconv.Itoa(int(e.GetStatusCode()))},
	}
	st, _ = st.WithDetails(detail)
	return st
}

func FromStatus(st *status.Status) (*Error, bool) {
	if st.Code() <= maxGRPCCode { // code from 0 to 17 are used in grpc std lib
		return nil, false
	}
	e := &Error{
		Code:       int32(st.Code()),
		Message:    st.Message(),
		Reason:     "",
		StatusCode: 0,
	}
loop:
	for _, detail := range st.Details() {
		switch val := detail.(type) {
		case *errdetails.ErrorInfo:
			e.Reason = val.Reason
			code, _ := strconv.Atoi(val.Metadata["status_code"])
			e.Code = int32(code)
			break loop
		}
	}
	return e, true
}

func FromError(err error) (*Error, bool) {
	var ee *Error
	if errors.As(err, &ee) {
		return ee, true
	}
	if st, ok := status.FromError(err); ok {
		return FromStatus(st)
	}
	return nil, false
}

func CodeFromError[T ~int32](err error) (T, bool) {
	if ee, ok := FromError(err); ok {
		return T(ee.GetCode()), true
	}
	return 0, false
}

func (e *Error) Unwrap() error {
	return e.Status().Err()
}

func New(code int32, message string, reason ...string) *Error {
	e := &Error{
		Code:    code,
		Message: message,
	}
	if len(reason) > 0 {
		e.Reason = reason[0]
	}
	return e
}
