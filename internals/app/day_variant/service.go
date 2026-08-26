package day_variant

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Create(req CreateRequest) (Response, error)
	ListByPlanIDFull(planID string) ([]Response, error)
	Update(req UpdateRequest) (Response, error)
	Delete(id string) error
}

type service struct {
	repo Repo
	log  *slog.Logger
}

// NewService builds the day_variant Service. log receives entries for
// mutating operations (Create/Update/Delete) - reads aren't logged.
func NewService(repo Repo, log *slog.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Create(req CreateRequest) (Response, error) {
	if req.PlanID == "" || req.Name == "" {
		return Response{}, ErrInvalidInput
	}

	dv := &DayVariant{
		PlanID: req.PlanID,
		Name:   req.Name,
		Order:  req.Order,
	}

	id, err := s.repo.Create(dv)
	if err != nil {
		s.log.Error("failed to create day variant", "plan_id", req.PlanID, "name", req.Name, "error", err)
		return Response{}, fmt.Errorf("create day variant: %w", err)
	}

	dv.ID = id
	s.log.Info("day variant created", "id", dv.ID, "plan_id", dv.PlanID, "name", dv.Name)
	return ToResponse(*dv), nil
}

func (s *service) ListByPlanIDFull(planID string) ([]Response, error) {
	dayVariants, err := s.repo.ListByPlanIDFull(planID)
	if err != nil {
		return nil, err
	}
	out := make([]Response, len(dayVariants))
	for i, dv := range dayVariants {
		out[i] = ToResponse(dv)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" || req.Name == "" {
		return Response{}, ErrInvalidInput
	}
	dv, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}
	dv.Name = req.Name
	dv.Order = req.Order
	if err := s.repo.Update(dv); err != nil {
		s.log.Error("failed to update day variant", "id", dv.ID, "error", err)
		return Response{}, fmt.Errorf("update day variant: %w", err)
	}
	s.log.Info("day variant updated", "id", dv.ID, "name", dv.Name)
	return ToResponse(*dv), nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete day variant", "id", id, "error", err)
		return err
	}
	s.log.Info("day variant deleted", "id", id)
	return nil
}
