package meal

import (
	"nutrition-planner/internals/app/food"
	"nutrition-planner/internals/app/supplement"
)

type Meal struct {
	ID           string `db:"id,primary_key,uuid"`
	DayVariantID string `db:"day_variant_id,notnull,references:day_variants(id),on_delete:cascade,index"`
	Name         string `db:"name,notnull"` // Name of the meal (e.g., "Breakfast, Lunch, Dinner, etc.")
	Note         string `db:"note"`         // Small note about the meal (e.g., "At 2-3 hours after the previous one")
	Order        int    `db:"order,default:0"`

	// Manually populated upon reading
	Foods       []food.Food             `db:"-"`
	Supplements []supplement.Supplement `db:"-"`
}

type CreateRequest struct {
	DayVariantID string `json:"dayVariantId"`
	Name         string `json:"name"`
	Note         string `json:"note"`
	Order        int    `json:"order"`
}

type UpdateRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Note  string `json:"note"`
	Order int    `json:"order"`
}

// Response — the data sent to the frontend. Modelul de domeniu (cu tag-uri
// `db`) nu ajunge niciodată direct în UI.
type Response struct {
	ID           string                `json:"id"`
	DayVariantID string                `json:"dayVariantId"`
	Name         string                `json:"name"`
	Note         string                `json:"note"`
	Order        int                   `json:"order"`
	Foods        []food.Response       `json:"foods,omitempty"`
	Supplements  []supplement.Response `json:"supplements,omitempty"`
}

// ToResponse e exportat pentru că day_variant.Response are nevoie să-l
// apeleze la construirea listei de mese dintr-o zi.
func ToResponse(m Meal) Response {
	foods := make([]food.Response, len(m.Foods))
	for i, f := range m.Foods {
		foods[i] = food.ToResponse(f)
	}
	supplements := make([]supplement.Response, len(m.Supplements))
	for i, s := range m.Supplements {
		supplements[i] = supplement.ToResponse(s)
	}
	return Response{
		ID:           m.ID,
		DayVariantID: m.DayVariantID,
		Name:         m.Name,
		Note:         m.Note,
		Order:        m.Order,
		Foods:        foods,
		Supplements:  supplements,
	}
}

func (Meal) TableName() string {
	return "meals"
}
