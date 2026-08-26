package supplement

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Create(req CreateRequest) (Response, error)
	ListByMealID(mealID string) ([]Response, error)
	Update(req UpdateRequest) (Response, error)
	Delete(id string) error
}

type service struct {
	repo Repo
	log  *slog.Logger
}

// NewService builds the supplement Service. log receives entries for
// mutating operations (Create/Update/Delete) - reads aren't logged.
func NewService(repo Repo, log *slog.Logger) Service { return &service{repo: repo, log: log} }

func (s *service) Create(req CreateRequest) (Response, error) {
	if req.MealID == "" || req.Supplement == "" {
		return Response{}, ErrInvalidInput
	}

	sup := &Supplement{
		MealID:     req.MealID,
		Supplement: req.Supplement,
		Quantity:   req.Quantity,
		Order:      req.Order,
	}

	id, err := s.repo.Create(sup)
	if err != nil {
		s.log.Error("failed to create supplement", "meal_id", req.MealID, "supplement", req.Supplement, "error", err)
		return Response{}, fmt.Errorf("create supplement: %w", err)
	}
	sup.ID = id
	s.log.Info("supplement created", "id", sup.ID, "meal_id", sup.MealID, "supplement", sup.Supplement)
	return ToResponse(*sup), nil
}

func (s *service) ListByMealID(mealID string) ([]Response, error) {
	supplements, err := s.repo.ListByMealID(mealID)
	if err != nil {
		return nil, err
	}
	out := make([]Response, len(supplements))
	for i, sup := range supplements {
		out[i] = ToResponse(sup)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" {
		return Response{}, ErrInvalidInput
	}
	sup, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}

	sup.Supplement = req.Supplement
	sup.Quantity = req.Quantity
	sup.Order = req.Order
	if err := s.repo.Update(sup); err != nil {
		s.log.Error("failed to update supplement", "id", sup.ID, "error", err)
		return Response{}, fmt.Errorf("update supplement: %w", err)
	}
	s.log.Info("supplement updated", "id", sup.ID, "supplement", sup.Supplement)
	return ToResponse(*sup), nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete supplement", "id", id, "error", err)
		return err
	}
	s.log.Info("supplement deleted", "id", id)
	return nil
}
