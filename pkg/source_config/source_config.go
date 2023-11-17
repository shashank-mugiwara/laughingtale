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

	sourceConfigsDto := type_configs.SourceConfigsDto{}
	if err := c.BodyParser(&sourceConfigsDto); err != nil {
		return err
	}

	errs := type_common.GetLaughingTaleValidator().Validate(sourceConfigsDto)
	if len(errs) != 0 {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errs, " and "),
		}
	}

	sourceConfigList := sourceConfigsDto.SourceConfig
	for _, scfg := range sourceConfigList {
		errs := type_common.GetLaughingTaleValidator().Validate(scfg)
		if len(errs) != 0 {
			return c.Status(fiber.ErrBadRequest.Code).JSON(&fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: strings.Join(errs, " and "),
			})
		}
	}

	sourceConfigs := type_configs.SourceConfigs{
		Identifier:   sourceConfigsDto.Identifier,
		SourceConfig: type_configs.JSONB{"configList": sourceConfigList},
	}
	if result := h.DB.Create(&sourceConfigs); result.Error != nil {
		h.Logger.Debug("Unable to create source loader config. Error is: ", result.Error.Error())
		dbErr := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    result.Error.Error(),
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(dbErr)
	}

	return c.Status(fiber.StatusCreated).JSON(sourceConfigs)
}
