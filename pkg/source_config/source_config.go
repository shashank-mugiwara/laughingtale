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

func (h handler) AddLoaderSourceConfig(c *fiber.Ctx) error {

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
		Type:         "loader",
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

func (h handler) GetLoaderSourceConfig(c *fiber.Ctx) error {
	h.Logger.Info("Got request for fetching loader source config with name: ", c.Params("name"))

	configName := c.Params("name")

	loaderScenarioConfig := type_configs.SourceConfigs{}
	result := h.DB.Find(&loaderScenarioConfig, "identifier = ?", configName)
	if result.Error != nil {
		h.Logger.Debug("Unable to fetch source loader config. Error is: ", result.Error.Error())
		dbErr := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    result.Error.Error(),
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(dbErr)
	}

	if result.RowsAffected == 0 {
		noRecFound := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    "No record exists with the given config name",
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(noRecFound)
	}

	return c.Status(fiber.StatusOK).JSON(loaderScenarioConfig)
}

func (h handler) DeleteLoaderSourceConfig(c *fiber.Ctx) error {
	h.Logger.Info("Got request for deleting loader source config with name: ", c.Params("name"))

	configName := c.Params("name")

	loaderScenarioConfig := type_configs.SourceConfigs{}
	result := h.DB.Find(&loaderScenarioConfig, "identifier = ?", configName)
	if result.Error != nil {
		h.Logger.Debug("Unable to fetch source loader config. Error is: ", result.Error.Error())
		dbErr := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    result.Error.Error(),
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(dbErr)
	}

	if result.RowsAffected == 0 {
		noRecFound := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    "No record exists with the given config name",
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(noRecFound)
	}

	result = h.DB.Delete(&loaderScenarioConfig, "identifier = ?", configName)
	if result.RowsAffected < 1 {
		noRecFound := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    "Config was not successfully deleted",
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(noRecFound)
	}

	return c.Status(fiber.StatusOK).JSON(type_common.Message{
		StatusCode: fiber.StatusOK,
		Message:    "Config successfully deleted",
	})
}

func (h handler) UpdateLoaderSourceConfig(c *fiber.Ctx) error {
	h.Logger.Info("Got request for deleting loader source config with name: ", c.Params("name"))

	configName := c.Params("name")

	loaderScenarioConfig := type_configs.SourceConfigs{}
	result := h.DB.Find(&loaderScenarioConfig, "identifier = ?", configName)
	if result.Error != nil {
		h.Logger.Debug("Unable to fetch source loader config. Error is: ", result.Error.Error())
		dbErr := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    result.Error.Error(),
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(dbErr)
	}

	if result.RowsAffected == 0 {
		noRecFound := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    "No record exists with the given config name",
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(noRecFound)
	}

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
		Type:         "loader",
	}
	if result := h.DB.Model(&sourceConfigs).Where("identifier = ?", configName).Updates(&sourceConfigs); result.Error != nil {
		h.Logger.Debug("Unable to create source loader config. Error is: ", result.Error.Error())
		dbErr := type_common.DatabaseError{
			StatusCode: fiber.ErrBadRequest.Code,
			Message:    result.Error.Error(),
		}

		return c.Status(fiber.ErrBadRequest.Code).JSON(dbErr)
	}

	return c.Status(fiber.StatusOK).JSON(sourceConfigs)
}
