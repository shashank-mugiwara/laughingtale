package type_connectors

type ConnectorConfigRequest struct {
	Name   string            `json:"name" validate:"required" `
	Config map[string]string `json:"config" validate:"required"`
}
