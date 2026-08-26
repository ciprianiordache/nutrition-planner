package day_variant

import "nutrition-planner/internals/app/meal"

type DayVariant struct {
	ID     string `db:"id,primary_key,uuid"`
	PlanID string `db:"plan_id,notnull,references:plans(id),on_delete:cascade,index"`
	Name   string `db:"name,notnull"` // Name of the day (e.g., "Working Day, Rest Day, etc.")
	Order  int    `db:"order,notnull,default:0"`

	// Manually populated upon reading
	Meals []meal.Meal `db:"-"`
}

type CreateRequest struct {
	PlanID string `json:"planId"`
	Name   string `json:"name"`
	Order  int    `json:"order"`
}

type UpdateRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// Response — the data sent to the frontend. Modelul de domeniu (cu tag-uri
// `db`) nu ajunge niciodată direct în UI.
type Response struct {
	ID     string          `json:"id"`
	PlanID string          `json:"planId"`
	Name   string          `json:"name"`
	Order  int             `json:"order"`
	Meals  []meal.Response `json:"meals,omitempty"`
}

// ToResponse e exportat pentru că plan.FullResponse are nevoie să-l apeleze
// la construirea listei de variante de zi ale unui plan.
func ToResponse(dv DayVariant) Response {
	meals := make([]meal.Response, len(dv.Meals))
	for i, m := range dv.Meals {
		meals[i] = meal.ToResponse(m)
	}
	return Response{
		ID:     dv.ID,
		PlanID: dv.PlanID,
		Name:   dv.Name,
		Order:  dv.Order,
		Meals:  meals,
	}
}

func (DayVariant) TableName() string {
	return "day_variants"
}
