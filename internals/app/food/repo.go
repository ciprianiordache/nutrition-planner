package food

import (
	"errors"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(f *Food) (string, error)
	GetByID(id string) (*Food, error)
	ListByMealID(mealID string) ([]Food, error)
	Update(f *Food) error
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

func (r *repo) Create(f *Food) (string, error) {
	return r.depot.Create(f)
}

func (r *repo) GetByID(id string) (*Food, error) {
	var f Food
	err := r.depot.ReadOne(Food{}, "id", id, &f)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *repo) ListByMealID(mealID string) ([]Food, error) {
	rows, err := r.db.Query(`
		SELECT id, meal_id, base_food, base_quantity, base_unit,
		       alternative_food, alternative_quantity, alternative_unit, "order"
		FROM foods
		WHERE meal_id = ?
		ORDER BY "order" ASC
	`, mealID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Food
	for rows.Next() {
		var f Food
		if err := rows.Scan(
			&f.ID, &f.MealID, &f.BaseFood, &f.BaseQuantity, &f.BaseUnit,
			&f.AlternativeFood, &f.AlternativeQuantity, &f.AlternativeUnit, &f.Order,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *repo) Update(f *Food) error {
	return r.depot.Update(f, "id", f.ID)
}

func (r *repo) Delete(id string) error {
	return r.depot.Delete(&Food{}, "id", id)
}

func (r *repo) DeleteByMealID(mealID string) error {
	_, err := r.db.Exec(`DELETE FROM foods WHERE meal_id = ?`, mealID)
	return err
}
