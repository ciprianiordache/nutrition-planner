package supplement

import (
	"errors"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(s *Supplement) (string, error)
	GetByID(id string) (*Supplement, error)
	ListByMealID(mealID string) ([]Supplement, error)
	Update(s *Supplement) error
	Delete(id string) error
	DeleteByMealID(mealID string) error
}

type repo struct {
	depot *crud.CRUD
	db    *database.Database
}

func NewRepo(db *database.Database) Repo {
	return &repo{
		depot: crud.New(db, crud.SQLite{}),
		db:    db,
	}
}

func (r *repo) Create(s *Supplement) (string, error) {
	return r.depot.Create(s)
}

func (r *repo) GetByID(id string) (*Supplement, error) {
	var s Supplement
	err := r.depot.ReadOne(Supplement{}, "id", id, &s)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *repo) ListByMealID(mealID string) ([]Supplement, error) {
	rows, err := r.db.Query(`
		SELECT id, meal_id, supplement, quantity, "order"
		FROM supplements
		WHERE meal_id = ?
		ORDER BY "order" ASC
	`, mealID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Supplement
	for rows.Next() {
		var s Supplement
		if err := rows.Scan(&s.ID, &s.MealID, &s.Supplement, &s.Quantity, &s.Order); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *repo) Update(s *Supplement) error {
	return r.depot.Update(s, "id", s.ID)
}

func (r *repo) Delete(id string) error {
	return r.depot.Delete(&Supplement{}, "id", id)
}

func (r *repo) DeleteByMealID(mealID string) error {
	_, err := r.db.Exec(`DELETE FROM supplements WHERE meal_id = ?`, mealID)
	return err
}
