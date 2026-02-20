package controller

import (
	"strconv"

	"github.com/Kittipoom-pan/autopart-service/internal/helper"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase"
	parterror "github.com/Kittipoom-pan/autopart-service/pkg/error"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PartController struct {
	usecase usecase.PartUsecase
	logger  zerolog.Logger
}

func NewPartController(usecase usecase.PartUsecase) *PartController {
	return &PartController{
		usecase: usecase,
		logger:  log.With().Str("component", "part_controller").Logger(),
	}
}

func (h *PartController) GetAllParts(c *fiber.Ctx) error {
	h.logger.Debug().Msg("Get all parts request")

	parts, err := h.usecase.GetAllParts(c.Context())
	if err != nil {
		return helper.RespondError(c, err)
	}
	if len(parts) == 0 {
		return helper.RespondSuccess(c, fiber.StatusOK, parts, "Parts not found")
	}

	return helper.RespondSuccess(c, fiber.StatusOK, parts, "")
}

func (h *PartController) GetPartByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	h.logger.Debug().Str("part_id", idStr).Msg("Get part request")

	if idStr == "" {
		h.logger.Warn().Msg("Id parameter is missing")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "id is required"}))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn().Msg("Invalid id format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid id format"}))
	}

	part, err := h.usecase.GetPartByID(c.Context(), id)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusOK, part, "")
}

func (h *PartController) CreatePart(c *fiber.Ctx) error {
	part := new(entitie.PartReq)
	if err := c.BodyParser(part); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	partID, err := h.usecase.CreatePart(c.Context(), part)
	if err != nil {
		return helper.RespondError(c, err)
	}

	return helper.RespondSuccess(c, fiber.StatusCreated, fiber.Map{"part_id": partID}, "Part created successfully")
}

func (h *PartController) UpdatePart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid part ID format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid part ID format"}))
	}

	partReq := new(entitie.PartReq)
	if err := c.BodyParser(partReq); err != nil {
		h.logger.Warn().Err(err).Bytes("raw_body", c.Body()).Msg("Failed to parse request body")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"body": "failed to parse request body"}))
	}

	if err := h.usecase.UpdatePart(c.Context(), id, partReq); err != nil {
		h.logger.Error().Err(err).Int("part_id", id).Msg("Failed to update part")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("part_id", id).Msg("Part updated successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Part updated successfully")
}

func (h *PartController) DeletePart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid part ID format")
		return helper.RespondError(c, parterror.InvalidRequestData(map[string]string{"id": "invalid part ID format"}))
	}

	if err := h.usecase.DeletePart(c.Context(), id); err != nil {
		h.logger.Error().Err(err).Int("part_id", id).Msg("Failed to delete part")
		return helper.RespondError(c, err)
	}

	h.logger.Info().Int("part_id", id).Msg("Part deleted successfully")
	return helper.RespondSuccess(c, fiber.StatusOK, nil, "Part deleted successfully")
}
