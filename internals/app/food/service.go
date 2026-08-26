package food

import (
	"database/sql"
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

// NewService builds the food Service. log receives entries for mutating
// operations (Create/Update/Delete) - reads aren't logged.
func NewService(repo Repo, log *slog.Logger) Service { return &service{repo: repo, log: log} }

func (s *service) Create(req CreateRequest) (Response, error) {
	if req.MealID == "" || req.BaseFood == "" || req.BaseQuantity <= 0 || req.BaseUnit == "" {
		return Response{}, ErrInvalidInput
	}

	f := &Food{
		MealID:              req.MealID,
		BaseFood:            req.BaseFood,
		BaseQuantity:        req.BaseQuantity,
		BaseUnit:            req.BaseUnit,
		AlternativeFood:     toNullString(req.AlternativeFood),
		AlternativeQuantity: toNullFloat(req.AlternativeQuantity),
		AlternativeUnit:     toNullString(req.AlternativeUnit),
		Order:               req.Order,
	}

	id, err := s.repo.Create(f)
	if err != nil {
		s.log.Error("failed to create food", "meal_id", req.MealID, "base_food", req.BaseFood, "error", err)
		return Response{}, fmt.Errorf("create food: %w", err)
	}
	f.ID = id
	s.log.Info("food created", "id", f.ID, "meal_id", f.MealID, "base_food", f.BaseFood)
	return ToResponse(*f), nil
}

func (s *service) ListByMealID(mealID string) ([]Response, error) {
	foods, err := s.repo.ListByMealID(mealID)
	if err != nil {
		return nil, err
	}
	out := make([]Response, len(foods))
	for i, f := range foods {
		out[i] = ToResponse(f)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" {
		return Response{}, ErrInvalidInput
	}
	f, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}
	f.BaseFood, f.BaseQuantity, f.BaseUnit, f.Order = req.BaseFood, req.BaseQuantity, req.BaseUnit, req.Order
	f.AlternativeFood = toNullString(req.AlternativeFood)
	f.AlternativeQuantity = toNullFloat(req.AlternativeQuantity)
	f.AlternativeUnit = toNullString(req.AlternativeUnit)

	if err := s.repo.Update(f); err != nil {
		s.log.Error("failed to update food", "id", f.ID, "error", err)
		return Response{}, fmt.Errorf("update food: %w", err)
	}
	s.log.Info("food updated", "id", f.ID, "base_food", f.BaseFood)
	return ToResponse(*f), nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete food", "id", id, "error", err)
		return err
	}
	s.log.Info("food deleted", "id", id)
	return nil
}

func toNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func toNullFloat(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}
