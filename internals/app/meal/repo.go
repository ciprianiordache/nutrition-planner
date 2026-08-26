package meal

import (
	"errors"
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/supplement"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(m *Meal) (string, error)
	GetByID(id string) (*Meal, error)
	ListByDayVariantID(dayVariantID string) ([]Meal, error)
	// ListByDayVariantIDFull — same as ListByDayVariantID, but also populates Foods/Supplements.
	ListByDayVariantIDFull(dayVariantID string) ([]Meal, error)
	Update(m *Meal) error
	Delete(id string) error
	DeleteByDayVariantID(dayVariantID string) error
}

type repo struct {
	depot     *crud.CRUD
	db        *database.Database
	foodRepo  food.Repo
	supplRepo supplement.Repo
}

func NewRepo(db *database.Database, foodRepo food.Repo, supplRepo supplement.Repo) Repo {
	return &repo{
		depot:     crud.New(db, crud.SQLite{}),
		db:        db,
		foodRepo:  foodRepo,
		supplRepo: supplRepo,
	}
}

func (r *repo) Create(m *Meal) (string, error) {
	return r.depot.Create(m)
}

func (r *repo) GetByID(id string) (*Meal, error) {
	var m Meal
	err := r.depot.ReadOne(Meal{}, "id", id, &m)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *repo) ListByDayVariantID(dayVariantID string) ([]Meal, error) {
	rows, err := r.db.Query(`
		SELECT id, day_variant_id, name, note, "order"
		FROM meals
		WHERE day_variant_id = ?
		ORDER BY "order" ASC
	`, dayVariantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Meal
	for rows.Next() {
		var m Meal
		if err := rows.Scan(&m.ID,
			&m.DayVariantID,
			&m.Name, &m.Note,
			&m.Order); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *repo) ListByDayVariantIDFull(dayVariantID string) ([]Meal, error) {
	meals, err := r.ListByDayVariantID(dayVariantID)
	if err != nil {
		return nil, err
	}
	for i := range meals {
		if err := r.attachNested(&meals[i]); err != nil {
			return nil, err
		}
	}
	return meals, nil
}

// attachNested populates Foods and Supplements for a single Meal (also used by GetByIDFull).
func (r *repo) attachNested(m *Meal) error {
	foods, err := r.foodRepo.ListByMealID(m.ID)
	if err != nil {
		return err
	}
	m.Foods = foods

	supplements, err := r.supplRepo.ListByMealID(m.ID)
	if err != nil {
		return err
	}
	m.Supplements = supplements

	return nil
}

func (r *repo) Update(m *Meal) error {
	return r.depot.Update(m, "id", m.ID)
}

func (r *repo) Delete(id string) error {
	// cascade is managed by the FK on_delete:cascade at DB level (SQLite),
	// so, foods/supplements are automatically deleted. Confirm that PRAGMA foreign_keys=ON is enabled.
	return r.depot.Delete(&Meal{}, "id", id)
}

func (r *repo) DeleteByDayVariantID(dayVariantID string) error {
	_, err := r.db.Exec(`DELETE FROM meals WHERE day_variant_id = ?`, dayVariantID)
	return err
}
