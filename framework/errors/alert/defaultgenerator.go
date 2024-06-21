package alert

var defaultGenerator = AlertErrorGenerator{}

// New returns an error object for the code, message.
func New(code int, reason, message string) *AlertError {
	return defaultGenerator.New(code, reason, message)
}

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(reason, message string) *AlertError {
	return defaultGenerator.Newf(400, reason, message)
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(reason, message string) *AlertError {
	return defaultGenerator.Newf(401, reason, message)
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(reason, message string) *AlertError {
	return defaultGenerator.Newf(403, reason, message)
}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(reason, message string) *AlertError {
	return defaultGenerator.Newf(404, reason, message)
}

// Conflict new Conflict error that is mapped to a 409 response.
func Conflict(reason, message string) *AlertError {
	return defaultGenerator.Newf(409, reason, message)
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func InternalServer(reason, message string) *AlertError {
	return defaultGenerator.Newf(500, reason, message)
}

// ServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ServiceUnavailable(reason, message string) *AlertError {
	return defaultGenerator.Newf(503, reason, message)
}

// GatewayTimeout new GatewayTimeout error that is mapped to a HTTP 504 response.
func GatewayTimeout(reason, message string) *AlertError {
	return defaultGenerator.Newf(504, reason, message)
}

// ClientClosed new ClientClosed error that is mapped to a HTTP 499 response.
func ClientClosed(reason, message string) *AlertError {
	return defaultGenerator.Newf(499, reason, message)
}
