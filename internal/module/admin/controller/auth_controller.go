package controller

import (
	"github.com/Kittipoom-pan/autopart-service/config"
	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/usecase"
	adminerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	usecase usecase.AuthUsecase
	logger  zerolog.Logger
	cfg     *config.Config
}

func NewAuthController(usecase usecase.AuthUsecase, cfg *config.Config) *AuthController {
	return &AuthController{
		usecase: usecase,
		cfg:     cfg,
		logger:  log.With().Str("component", "auth_controller").Logger(),
	}
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	request := new(entitie.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, adminerror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	response, err := h.usecase.Login(c.Context(), request)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, response, "")
}
