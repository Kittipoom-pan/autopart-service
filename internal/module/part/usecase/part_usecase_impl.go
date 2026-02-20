package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/part/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/repository"
	"github.com/Kittipoom-pan/autopart-service/internal/module/part/usecase/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type partUsecase struct {
	repo   repository.PartRepository
	logger zerolog.Logger
}

func NewPartUsecase(repo repository.PartRepository) PartUsecase {
	return &partUsecase{
		repo:   repo,
		logger: log.With().Str("component", "part_usecase").Logger(),
	}
}

func (u *partUsecase) GetPartByID(ctx context.Context, id int) (*entitie.PartRes, error) {
	u.logger.Info().Int("part_id", id).Msg("GetPartByID started")

	part, err := u.repo.GetPartByID(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("part_id", id).Msg("Failed to get part from repository")
		return nil, err
	}

	u.logger.Info().Int("part_id", id).Msg("GetPartByID completed successfully")
	return part, nil
}

func (u *partUsecase) GetAllParts(ctx context.Context) ([]*entitie.PartRes, error) {
	parts, err := u.repo.GetAllParts(ctx)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to get all parts from repository")
		return nil, err
	}
	return parts, nil
}

func (u *partUsecase) CreatePart(ctx context.Context, part *entitie.PartReq) (int64, error) {
	partID, err := u.repo.CreatePart(ctx, part)
	if err != nil {
		u.logger.Error().Err(err).Msg("Failed to create part")
		return 0, err
	}
	u.logger.Info().Int64("part_id", partID).Msg("Part created successfully")
	return partID, nil
}

func (u *partUsecase) UpdatePart(ctx context.Context, id int, partReq *entitie.PartReq) error {
	u.logger.Info().Int("part_id", id).Msg("UpdatePart started")

	if err := validation.ValidatePartRequest(partReq, true); err != nil {
		u.logger.Warn().Err(err).Msg("Part request validation failed")
		return err
	}

	_, err := u.repo.GetPartByID(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("part_id", id).Msg("Part not found or failed to get from repository")
		return err
	}

	err = u.repo.UpdatePart(ctx, id, partReq)
	if err != nil {
		u.logger.Error().Err(err).Int("part_id", id).Msg("Failed to update part in repository")
		return err
	}

	u.logger.Info().Int("part_id", id).Msg("UpdatePart completed successfully")
	return nil
}

func (u *partUsecase) DeletePart(ctx context.Context, id int) error {
	u.logger.Info().Int("part_id", id).Msg("DeletePart started")

	err := u.repo.DeletePart(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Int("part_id", id).Msg("Failed to delete part in repository")
		return err
	}

	u.logger.Info().Int("part_id", id).Msg("DeletePart completed successfully")
	return nil
}
