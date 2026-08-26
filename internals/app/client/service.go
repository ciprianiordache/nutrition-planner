package client

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Create(req CreateRequest) (Response, error)
	GetByID(id string) (Response, error)
	List(offset, limit int) ([]Response, error)
	Update(req UpdateRequest) (Response, error)
	Delete(id string) error
}

type service struct {
	repo Repo
	log  *slog.Logger
}

// NewService builds the client Service. log receives entries for mutating
// operations (Create/Update/Delete) - reads (GetByID/List) aren't logged,
// they're too frequent and carry no operational signal on their own.
func NewService(repo Repo, log *slog.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) Create(req CreateRequest) (Response, error) {
	if err := validate(req.Name, req.Weight, req.Height); err != nil {
		return Response{}, err
	}

	c := &Client{
		Name:   req.Name,
		Weight: req.Weight,
		Height: req.Height,
		Target: req.Target,
	}

	id, err := s.repo.Create(c)
	if err != nil {
		s.log.Error("failed to create client", "name", c.Name, "error", err)
		return Response{}, fmt.Errorf("create client: %w", err)
	}
	c.ID = id
	s.log.Info("client created", "id", c.ID, "name", c.Name)
	return toResponse(*c), nil
}

func (s *service) GetByID(id string) (Response, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return Response{}, err
	}
	return toResponse(*c), nil
}

func (s *service) List(offset, limit int) ([]Response, error) {
	clients, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, err
	}

	out := make([]Response, len(clients))
	for i, c := range clients {
		out[i] = toResponse(c)
	}
	return out, nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	if req.ID == "" {
		return Response{}, ErrInvalidInput
	}
	if err := validate(req.Name, req.Weight, req.Height); err != nil {
		return Response{}, err
	}

	c, err := s.repo.GetByID(req.ID)
	if err != nil {
		return Response{}, err
	}
	c.Name = req.Name
	c.Weight = req.Weight
	c.Height = req.Height
	c.Target = req.Target

	if err := s.repo.Update(c); err != nil {
		s.log.Error("failed to update client", "id", c.ID, "error", err)
		return Response{}, fmt.Errorf("update client: %w", err)
	}
	s.log.Info("client updated", "id", c.ID, "name", c.Name)
	return toResponse(*c), nil
}

func (s *service) Delete(id string) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("failed to delete client", "id", id, "error", err)
		return err
	}
	s.log.Info("client deleted", "id", id)
	return nil
}

func validate(name string, weight, height float64) error {
	if name == "" || weight <= 0 || height <= 0 {
		return ErrInvalidInput
	}
	return nil
}
