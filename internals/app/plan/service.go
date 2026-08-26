package plan

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Create(req CreateRequest) (Response, error)
	// GetByIDFull returns the full tree (Plan + DayVariants -> Meals ->
	// Foods/Supplements), used when editing the plan and generating the PDF.
	GetByIDFull(id string) (FullResponse, error)
	ListByClientID(clientID string) ([]Response, error)
	Update(req UpdateRequest) (Response, error)
	Delete(id string) error
}

type service struct {
	repo Repo
	log  *slog.Logger
}

// NewService builds the plan Service. log receives entries for mutating
// operations (Create/Update/Delete) - reads aren't logged.
func NewService(repo Repo, log *slog.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) Create(req CreateRequest) (Response, error) {
	if req.ClientID == "" || req.Title == "" {
		return Response{}, ErrInvalidInput
	}

	p := &Plan{
		ClientID:               req.ClientID,
		Title:                  req.Title,
		GeneralRecommendations: req.GeneralRecommendations,
		Footer:                 req.Footer,
	}
	id, err := s.repo.Create(p)
	if err != nil {
		s.log.Error("failed to create plan", "client_id", req.ClientID, "title", req.Title, "error", err)
		return Response{}, fmt.Errorf("create plan: %w", err)
	}
	p.ID = id
	s.log.Info("plan created", "id", p.ID, "client_id", p.ClientID, "title", p.Title)
	return toResponse(*p), nil
}

func (s *service) GetByIDFull(id string) (FullResponse, error) {
	p, err := s.repo.GetByIDFull(id)
	if err != nil {
		return FullResponse{}, err
	}
	return toFullResponse(*p), nil
}

func (s *service) ListByClientID(clientID string) ([]Response, error) {
	plans, err := s.repo.ListByClientID(clientID)
	if err != nil {
		return nil, err
	}

	out := make([]Response, len(plans))
	for i, p := range plans {
		out[i] = toResponse(p)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" || req.Title == "" {
		return Response{}, ErrInvalidInput
	}

	p, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}
	p.Title = req.Title
	p.GeneralRecommendations = req.GeneralRecommendations
	p.Footer = req.Footer

	if err := s.repo.Update(p); err != nil {
		s.log.Error("failed to update plan", "id", p.ID, "error", err)
		return Response{}, fmt.Errorf("update plan: %w", err)
	}
	s.log.Info("plan updated", "id", p.ID, "title", p.Title)
	return toResponse(*p), nil
}

func (s *service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete plan", "id", id, "error", err)
		return err
	}
	s.log.Info("plan deleted", "id", id)
	return nil
}
