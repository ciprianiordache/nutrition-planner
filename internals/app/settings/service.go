package settings

import (
	"fmt"
	"log/slog"
)

type Service interface {
	Get() (Response, error)
	Update(req UpdateRequest) (Response, error)
}

type service struct {
	repo Repo
	log  *slog.Logger
}

func NewService(repo Repo, log *slog.Logger) Service { return &service{repo: repo, log: log} }

func (s *service) Get() (Response, error) {
	settings, err := s.repo.Get()
	if err != nil {
		return Response{}, err
	}
	return toResponse(*settings), nil
}

func (s *service) Update(req UpdateRequest) (Response, error) {
	current, err := s.repo.Get()
	if err != nil {
		return Response{}, err
	}

	applyUpdate(current, req)

	if err := s.repo.Update(current); err != nil {
		s.log.Error("failed to update settings", "error", err)
		return Response{}, fmt.Errorf("update settings: %w", err)
	}
	s.log.Info("settings updated")
	return toResponse(*current), nil
}
