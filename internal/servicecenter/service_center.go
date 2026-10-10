package servicecenter

import "vyshkat/internal/domain"

// ServiceCenter performs the technical inspection before transport enters the active fleet.
type ServiceCenter struct{}

func New() *ServiceCenter {
	return &ServiceCenter{}
}

// Inspect accepts operational transport and rejects broken transport.
// The concrete inspection rule is deliberately small: the homework does not prescribe a richer one.
func (s *ServiceCenter) Inspect(transport domain.Transport) bool {
	return transport != nil && transport.IsOperational()
}
