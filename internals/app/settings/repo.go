package settings

import (
	"errors"

	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

// singletonID is the fixed primary key of the app's one-and-only settings
// row - there is no per-user or per-profile concept in this app, so
// there's no reason to let more than one row ever exist.
const singletonID = "app-settings"

type Repo interface {
	// Get returns the settings row, creating it with all-empty (i.e. "use
	// config.yaml defaults") values on first call.
	Get() (*Settings, error)
	Update(s *Settings) error
}

type repo struct {
	depot *crud.CRUD
}

func NewRepo(db *database.Database) Repo { return &repo{depot: crud.New(db, crud.SQLite{})} }

func (r *repo) Get() (*Settings, error) {
	var s Settings
	err := r.depot.ReadOne(Settings{}, "id", singletonID, &s)
	if errors.Is(err, crud.ErrNotFound) {
		s = Settings{ID: singletonID}
		if _, createErr := r.depot.Create(&s); createErr != nil {
			return nil, createErr
		}
		return &s, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repo) Update(s *Settings) error {
	s.ID = singletonID // never let a caller-supplied ID create a second row
	return r.depot.Update(s, "id", singletonID)
}
