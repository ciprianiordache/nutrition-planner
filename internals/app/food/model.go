package food

import "database/sql"

type Food struct {
	ID     string `db:"id,primary_key,uuid"`
	MealID string `db:"meal_id,notnull,references:meals(id),on_delete:cascade,index"`

	BaseFood     string  `db:"base_food,notnull"`     // Base food (e.g., "Chicken breast, Rice, Broccoli, etc.")
	BaseQuantity float64 `db:"base_quantity,notnull"` // Base Food quantity (e.g., "100, 200, 300, etc.")
	BaseUnit     string  `db:"base_unit,notnull"`     // Base Food unit (e.g., "g, ml, etc.")

	// The alternative food item is optional (not every row has an alternative in the example, e.g., "Green beans")
	AlternativeFood     sql.NullString  `db:"alternative_food"`
	AlternativeQuantity sql.NullFloat64 `db:"alternative_quantity"`
	AlternativeUnit     sql.NullString  `db:"alternative_unit"`

	Order int `db:"order,default:0"`
}

type CreateRequest struct {
	MealID              string   `json:"mealId"`
	BaseFood            string   `json:"baseFood"`
	BaseQuantity        float64  `json:"baseQuantity"`
	BaseUnit            string   `json:"baseUnit"`
	AlternativeFood     *string  `json:"alternativeFood"`
	AlternativeQuantity *float64 `json:"alternativeQuantity"`
	AlternativeUnit     *string  `json:"alternativeUnit"`
	Order               int      `json:"order"`
}

type UpdateRequest struct {
	ID                  string   `json:"id"`
	BaseFood            string   `json:"baseFood"`
	BaseQuantity        float64  `json:"baseQuantity"`
	BaseUnit            string   `json:"baseUnit"`
	AlternativeFood     *string  `json:"alternativeFood"`
	AlternativeQuantity *float64 `json:"alternativeQuantity"`
	AlternativeUnit     *string  `json:"alternativeUnit"`
	Order               int      `json:"order"`
}

// Response — the data sent to the frontend. Modelul de domeniu (cu tag-uri
// `db` și sql.Null*) nu ajunge niciodată direct în UI.
type Response struct {
	ID                  string   `json:"id"`
	MealID              string   `json:"mealId"`
	BaseFood            string   `json:"baseFood"`
	BaseQuantity        float64  `json:"baseQuantity"`
	BaseUnit            string   `json:"baseUnit"`
	AlternativeFood     *string  `json:"alternativeFood,omitempty"`
	AlternativeQuantity *float64 `json:"alternativeQuantity,omitempty"`
	AlternativeUnit     *string  `json:"alternativeUnit,omitempty"`
	Order               int      `json:"order"`
}

// ToResponse e exportat (nu doar folosit intern) pentru că meal.Response
// are nevoie să-l apeleze la construirea listei de foods dintr-o masă.
func ToResponse(f Food) Response {
	r := Response{
		ID:           f.ID,
		MealID:       f.MealID,
		BaseFood:     f.BaseFood,
		BaseQuantity: f.BaseQuantity,
		BaseUnit:     f.BaseUnit,
		Order:        f.Order,
	}
	if f.AlternativeFood.Valid {
		r.AlternativeFood = &f.AlternativeFood.String
	}
	if f.AlternativeQuantity.Valid {
		r.AlternativeQuantity = &f.AlternativeQuantity.Float64
	}
	if f.AlternativeUnit.Valid {
		r.AlternativeUnit = &f.AlternativeUnit.String
	}
	return r
}

func (Food) TableName() string {
	return "foods"
}
