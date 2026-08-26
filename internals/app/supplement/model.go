package supplement

type Supplement struct {
	ID         string  `db:"id,primary_key,uuid"`
	MealID     string  `db:"meal_id,notnull,references:meals(id),on_delete:cascade,index"`
	Supplement string  `db:"supplement"` // Supplement (e.g., "multivitamins, omega-3, etc.")
	Quantity   float64 `db:"quantity"`   // Supplement quantity (e.g., "1, 2, 3, etc.")
	Order      int     `db:"order,default:0"`
}

type CreateRequest struct {
	MealID     string  `json:"mealId"`
	Supplement string  `json:"supplement"`
	Quantity   float64 `json:"quantity"`
	Order      int     `json:"order"`
}

type UpdateRequest struct {
	ID         string  `json:"id"`
	Supplement string  `json:"supplement"`
	Quantity   float64 `json:"quantity"`
	Order      int     `json:"order"`
}

// Response — the data sent to the frontend. Modelul de domeniu (cu tag-uri
// `db`) nu ajunge niciodată direct în UI.
type Response struct {
	ID         string  `json:"id"`
	MealID     string  `json:"mealId"`
	Supplement string  `json:"supplement"`
	Quantity   float64 `json:"quantity"`
	Order      int     `json:"order"`
}

// ToResponse e exportat pentru că meal.Response are nevoie să-l apeleze la
// construirea listei de suplimente dintr-o masă.
func ToResponse(s Supplement) Response {
	return Response{
		ID:         s.ID,
		MealID:     s.MealID,
		Supplement: s.Supplement,
		Quantity:   s.Quantity,
		Order:      s.Order,
	}
}

func (Supplement) TableName() string {
	return "supplements"
}
