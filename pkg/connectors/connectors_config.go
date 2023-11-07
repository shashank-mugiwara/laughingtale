package connectors

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_common"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_connectors"
)

type handler struct {
	Logger     Logger
	HTTPClient *resty.Client
}

func NewHandler(logger Logger, http_client *resty.Client) Handler {
	return &handler{
		Logger:     logger,
		HTTPClient: http_client,
	}
}

func (h *handler) AddConnectorConfig(c *fiber.Ctx) error {
	connector_request_payload := type_connectors.ConnectorConfigRequest{}

	if err := c.BodyParser(&connector_request_payload); err != nil {
		return err
	}

	if errs := type_common.GetLaughingTaleValidator().Validate(connector_request_payload); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)
		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	return c.JSON(connector_request_payload)
}
