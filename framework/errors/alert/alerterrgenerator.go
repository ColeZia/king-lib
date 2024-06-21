package alert

import (
	"fmt"

	"gl.king.im/king-lib/framework/alerting"
)

type AlertErrorGenerator struct {
	Alerting *alerting.Alerting
	AlertFun alerting.AlertFun
}

// New returns an error object for the code, message.
func (aeg *AlertErrorGenerator) New(code int, reason, message string) *AlertError {
	return &AlertError{
		Code:     int32(code),
		Message:  message,
		Reason:   reason,
		Alerting: aeg.Alerting,
		AlertFun: aeg.AlertFun,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func (aeg *AlertErrorGenerator) Newf(code int, reason, format string, a ...interface{}) *AlertError {
	return aeg.New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func (aeg *AlertErrorGenerator) Errorf(code int, reason, format string, a ...interface{}) error {
	return aeg.New(code, reason, fmt.Sprintf(format, a...))
}

// BadRequest new BadRequest error that is mapped to a 400 response.
func (aeg *AlertErrorGenerator) BadRequest(reason, message string) *AlertError {
	return aeg.Newf(400, reason, message)
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func (aeg *AlertErrorGenerator) Unauthorized(reason, message string) *AlertError {
	return aeg.Newf(401, reason, message)
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func (aeg *AlertErrorGenerator) Forbidden(reason, message string) *AlertError {
	return aeg.Newf(403, reason, message)
}

// NotFound new NotFound error that is mapped to a 404 response.
func (aeg *AlertErrorGenerator) NotFound(reason, message string) *AlertError {
	return aeg.Newf(404, reason, message)
}

// Conflict new Conflict error that is mapped to a 409 response.
func (aeg *AlertErrorGenerator) Conflict(reason, message string) *AlertError {
	return aeg.Newf(409, reason, message)
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func (aeg *AlertErrorGenerator) InternalServer(reason, message string) *AlertError {
	return aeg.Newf(500, reason, message)
}

// ServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func (aeg *AlertErrorGenerator) ServiceUnavailable(reason, message string) *AlertError {
	return aeg.Newf(503, reason, message)
}

// GatewayTimeout new GatewayTimeout error that is mapped to a HTTP 504 response.
func (aeg *AlertErrorGenerator) GatewayTimeout(reason, message string) *AlertError {
	return aeg.Newf(504, reason, message)
}

// ClientClosed new ClientClosed error that is mapped to a HTTP 499 response.
func (aeg *AlertErrorGenerator) ClientClosed(reason, message string) *AlertError {
	return aeg.Newf(499, reason, message)
}
