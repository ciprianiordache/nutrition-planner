package day_variant

import (
	"errors"
	"nutrition-planner/internals/app/meal"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(dv *DayVariant) (string, error)
	GetByID(id string) (*DayVariant, error)
	ListByPlanID(planID string) ([]DayVariant, error)
	ListByPlanIDFull(planID string) ([]DayVariant, error)
	Update(dv *DayVariant) error
	Delete(id string) error
	DeleteByPlanID(planID string) error
	// Duplicate creates a copy of source (name/order, plus every meal it
	// has - each with its own foods/supplements) under planID, returning
	// the new day variant's ID. source is expected to already have Meals
	// populated (e.g. from ListByPlanIDFull) - Duplicate doesn't re-query
	// them.
	Duplicate(planID string, source DayVariant) (string, error)
}

type repo struct {
	depot    *crud.CRUD
	db       *database.Database
	mealRepo meal.Repo
}

func NewRepo(db *database.Database, mealRepo meal.Repo) Repo {
	return &repo{
		depot:    crud.New(db, crud.SQLite{}),
		db:       db,
		mealRepo: mealRepo,
	}
}

func (r *repo) Create(dv *DayVariant) (string, error) {
	return r.depot.Create(dv)
}

func (r *repo) GetByID(id string) (*DayVariant, error) {
	var dv DayVariant
	err := r.depot.ReadOne(DayVariant{}, "id", id, &dv)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &dv, nil
}

func (r *repo) ListByPlanID(planID string) ([]DayVariant, error) {
	rows, err := r.db.Query(`
		SELECT id, plan_id, name, "order"
		FROM day_variants
		WHERE plan_id = ?
		ORDER BY "order" ASC
	`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DayVariant
	for rows.Next() {
		var dv DayVariant
		if err := rows.Scan(&dv.ID, &dv.PlanID, &dv.Name, &dv.Order); err != nil {
			return nil, err
		}
		out = append(out, dv)
	}
	return out, rows.Err()

}

func (r *repo) ListByPlanIDFull(planID string) ([]DayVariant, error) {
	dayVariants, err := r.ListByPlanID(planID)
	if err != nil {
		return nil, err
	}

	for i := range dayVariants {
		meals, err := r.mealRepo.ListByDayVariantIDFull(dayVariants[i].ID)
		if err != nil {
			return nil, err
		}
		dayVariants[i].Meals = meals
	}
	return dayVariants, nil
}

func (r *repo) Update(dv *DayVariant) error {
	return r.depot.Update(dv, "id", dv.ID)
}

func (r *repo) Delete(id string) error {
	return r.depot.Delete(&DayVariant{}, "id", id)
}

func (r *repo) DeleteByPlanID(planID string) error {
	_, err := r.db.Exec(`DELETE FROM day_variants WHERE plan_id = ?`, planID)
	return err
}

// Duplicate creates a new day variant under planID with the same name/
// order as source, then duplicates every one of source's meals onto it
// (each with fresh IDs, recursively duplicating their own foods/
// supplements too - see meal.Repo.Duplicate). source.Meals must already
// be populated (see ListByPlanIDFull).
func (r *repo) Duplicate(planID string, source DayVariant) (string, error) {
	newDayVariant := &DayVariant{
		PlanID: planID,
		Name:   source.Name,
		Order:  source.Order,
	}
	newID, err := r.depot.Create(newDayVariant)
	if err != nil {
		return "", err
	}

	for _, m := range source.Meals {
		if _, err := r.mealRepo.Duplicate(newID, m); err != nil {
			return "", err
		}
	}

	return newID, nil
}
