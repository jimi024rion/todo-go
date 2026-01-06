package health

// Handler is a collection of health check handlers.
type Handler struct {
	Check *HealthCheckHandler
}

// NewHandler creates a new health check handler collection.
func NewHandler(healthCheckHandler *HealthCheckHandler) *Handler {
	return &Handler{
		Check: healthCheckHandler,
	}
}
