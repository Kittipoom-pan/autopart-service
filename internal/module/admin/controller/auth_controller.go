package controller

import (
	"github.com/Kittipoom-pan/autopart-service/config"
	_ "github.com/Kittipoom-pan/autopart-service/internal/common"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthController struct {
	usecase usecase.AuthUsecase
	logger  zerolog.Logger
	cfg     *config.Config
}

func NewAuthController(usecase usecase.AuthUsecase, cfg *config.Config, logger zerolog.Logger) *AuthController {
	return &AuthController{
		usecase: usecase,
		cfg:     cfg,
		logger:  logger,
	}
}

// Login godoc
// @Summary      Admin login
// @Description  เข้าสู่ระบบสำหรับผู้ดูแลระบบ
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body  body      entity.LoginRequest  true  "Login credentials"
// @Success      200  {object}  common.BaseResponse[entity.LoginResponse]
// @Failure      400  {object}  common.BaseErrorResponse
// @Failure      401  {object}  common.BaseErrorResponse
// @Failure      500  {object}  common.BaseErrorResponse
// @Router       /v1/admin/login [post]
func (h *AuthController) Login(c *fiber.Ctx) error {
	request := new(entity.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	response, err := h.usecase.Login(c.UserContext(), request)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, response, "")
}
