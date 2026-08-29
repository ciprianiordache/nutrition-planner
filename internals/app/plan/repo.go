package plan

import (
	"errors"
	"nutrition-planner/internals/app/day_variant"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(p *Plan) (string, error)
	GetByID(id string) (*Plan, error)
	GetByIDFull(id string) (*Plan, error)
	ListByClientID(clientID string) ([]Plan, error)
	Update(p *Plan) error
	Delete(id string) error
	// Duplicate creates a full copy of source (title, general
	// recommendations, footer, every day variant/meal/food/supplement it
	// has) for targetClientID (which may be a *different* client than
	// source.ClientID), returning the new plan's ID. source is expected to
	// come from GetByIDFull (its DayVariants must already be populated) -
	// Duplicate doesn't re-query them.
	Duplicate(source *Plan, targetClientID, newTitle string) (string, error)
}

type repo struct {
	depot      *crud.CRUD
	db         *database.Database
	dayVarRepo day_variant.Repo
}

func NewRepo(db *database.Database, dayVarRepo day_variant.Repo) Repo {
	return &repo{
		depot:      crud.New(db, crud.SQLite{}),
		db:         db,
		dayVarRepo: dayVarRepo,
	}
}

func (r *repo) Create(p *Plan) (string, error) {
	return r.depot.Create(p)
}

func (r *repo) GetByID(id string) (*Plan, error) {
	var p Plan
	err := r.depot.ReadOne(Plan{}, "id", id, &p)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetByIDFull — the complete plan, including day_variants -> meals -> foods/supplements.
// Use this for the plan editing screen and when generating the PDF.
func (r *repo) GetByIDFull(id string) (*Plan, error) {
	p, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	dayVariants, err := r.dayVarRepo.ListByPlanIDFull(p.ID)
	if err != nil {
		return nil, err
	}

	p.DayVariants = dayVariants
	return p, err
}

func (r *repo) ListByClientID(clientID string) ([]Plan, error) {
	rows, err := r.db.Query(`
		SELECT id,
			client_id,
			title, 
			general_recommendations, 
			footer, 
			created_at, 
			updated_at
		FROM plans
		WHERE client_id = ?
		ORDER BY created_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(
			&p.ID,
			&p.ClientID,
			&p.Title,
			&p.GeneralRecommendations,
			&p.Footer,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *repo) Update(p *Plan) error {
	return r.depot.Update(p, "id", p.ID)
}

func (r *repo) Delete(id string) error {
	return r.depot.Delete(&Plan{}, "id", id)
}

// Duplicate creates a new plan for targetClientID (which may be a
// *different* client than source.ClientID) titled newTitle, with the same
// general recommendations/footer, then duplicates every one of source's
// day variants onto it (each recursively duplicating its own meals, and
// each meal its own foods/supplements - see day_variant.Repo.Duplicate and
// meal.Repo.Duplicate). source.DayVariants must already be populated (see
// GetByIDFull).
func (r *repo) Duplicate(source *Plan, targetClientID, newTitle string) (string, error) {
	newPlan := &Plan{
		ClientID:               targetClientID,
		Title:                  newTitle,
		GeneralRecommendations: source.GeneralRecommendations,
		Footer:                 source.Footer,
	}
	newID, err := r.depot.Create(newPlan)
	if err != nil {
		return "", err
	}

	for _, dv := range source.DayVariants {
		if _, err := r.dayVarRepo.Duplicate(newID, dv); err != nil {
			return "", err
		}
	}

	return newID, nil
}
