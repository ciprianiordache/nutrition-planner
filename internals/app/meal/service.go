package meal

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Create(req CreateRequest) (Response, error)
	ListByDayVariantIDFull(dayVariantID string) ([]Response, error)
	Update(req UpdateRequest) (Response, error)
	Delete(id string) error
}

type service struct {
	repo Repo
	log  *slog.Logger
}

// NewService builds the meal Service. log receives entries for mutating
// operations (Create/Update/Delete) - reads aren't logged.
func NewService(repo Repo, log *slog.Logger) Service { return &service{repo: repo, log: log} }

func (s *service) Create(req CreateRequest) (Response, error) {
	if req.DayVariantID == "" || req.Name == "" {
		return Response{}, ErrInvalidInput
	}
	m := &Meal{DayVariantID: req.DayVariantID, Name: req.Name, Note: req.Note, Order: req.Order}
	id, err := s.repo.Create(m)
	if err != nil {
		s.log.Error("failed to create meal", "day_variant_id", req.DayVariantID, "name", req.Name, "error", err)
		return Response{}, fmt.Errorf("create meal: %w", err)
	}
	m.ID = id
	s.log.Info("meal created", "id", m.ID, "day_variant_id", m.DayVariantID, "name", m.Name)
	return ToResponse(*m), nil
}

func (s *service) ListByDayVariantIDFull(dayVariantID string) ([]Response, error) {
	meals, err := s.repo.ListByDayVariantIDFull(dayVariantID)
	if err != nil {
		return nil, err
	}
	out := make([]Response, len(meals))
	for i, m := range meals {
		out[i] = ToResponse(m)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" || req.Name == "" {
		return Response{}, ErrInvalidInput
	}
	m, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}
	m.Name, m.Note, m.Order = req.Name, req.Note, req.Order
	if err := s.repo.Update(m); err != nil {
		s.log.Error("failed to update meal", "id", m.ID, "error", err)
		return Response{}, fmt.Errorf("update meal: %w", err)
	}
	s.log.Info("meal updated", "id", m.ID, "name", m.Name)
	return ToResponse(*m), nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete meal", "id", id, "error", err)
		return err
	}
	s.log.Info("meal deleted", "id", id)
	return nil
}
