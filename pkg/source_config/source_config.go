package sourceconfig

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_common"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_configs"
	"gorm.io/gorm"
)

type handler struct {
	Logger Logger
	DB     *gorm.DB
}

func NewHandler(logger Logger, db *gorm.DB) Handler {
	return &handler{
		Logger: logger,
		DB:     db,
	}
}

func (h handler) AddSourceConfig(c *fiber.Ctx) error {
	h.Logger.Info("Got request for adding loader source config with name: ", c.Params("name"))

	sourceConfigs := type_configs.SourceConfigs{}
	if err := c.BodyParser(&sourceConfigs); err != nil {
		return err
	}

	errs := type_common.GetLaughingTaleValidator().Validate(sourceConfigs)
	if len(errs) != 0 {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errs, " and "),
		}
	}

	return nil
}
