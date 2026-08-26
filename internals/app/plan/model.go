package plan

import (
	"nutrition-planner/internals/app/day_variant"
	"time"
)

type Plan struct {
	ID                     string    `db:"id,primary_key,uuid"`
	ClientID               string    `db:"client_id,notnull,references:clients(id),on_delete:cascade,index"`
	Title                  string    `db:"title,notnull"`                                                                             // Title of the plan ("e.g., "Plan alimentar Client(Name)")
	GeneralRecommendations string    `db:"general_recommendations"`                                                                   // General recommendations (e.g, "Drink 3 liters of water, avoid fried foods, sleep at least 7 hours per day, etc.")
	Footer                 string    `db:"footer,default:Discount Code: VITABOLIC, cod MihailVTB - 10% reducere + transport gratuit"` // Footer of the plan (e.g., "You can add here any additional information or instruction for the client. " default: "Discount Code: VITABOLIC, cod MihailVTB - 10% reducere + transport gratuit")
	CreatedAt              time.Time `db:"created_at,oncreate"`
	UpdatedAt              time.Time `db:"updated_at,onwrite"`

	// It's not a column – crud-depot/schema-builder doesn't perform joins (they aren't ORMs).
	// It is populated manually in the repository after reading from day_variants.
	DayVariants []day_variant.DayVariant `db:"-"`
}

type CreateRequest struct {
	ClientID               string `json:"clientId"`
	Title                  string `json:"title"`
	GeneralRecommendations string `json:"generalRecommendations"`
	Footer                 string `json:"footer"`
}

type UpdateRequest struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	GeneralRecommendations string `json:"generalRecommendations"`
	Footer                 string `json:"footer"`
}

// Response — the data sent to the frontend pentru Create/List/Update.
// Modelul de domeniu (cu tag-uri `db`) nu ajunge niciodată direct în UI.
type Response struct {
	ID                     string    `json:"id"`
	ClientID               string    `json:"clientId"`
	Title                  string    `json:"title"`
	GeneralRecommendations string    `json:"generalRecommendations"`
	Footer                 string    `json:"footer"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// FullResponse — folosit doar de GetPlanFull, include arborele complet
// (day variants -> meals -> foods/supplements), pentru editorul de plan și
// generarea PDF-ului.
type FullResponse struct {
	Response
	DayVariants []day_variant.Response `json:"dayVariants,omitempty"`
}

func toResponse(p Plan) Response {
	return Response{
		ID:                     p.ID,
		ClientID:               p.ClientID,
		Title:                  p.Title,
		GeneralRecommendations: p.GeneralRecommendations,
		Footer:                 p.Footer,
		CreatedAt:              p.CreatedAt,
		UpdatedAt:              p.UpdatedAt,
	}
}

func toFullResponse(p Plan) FullResponse {
	dayVariants := make([]day_variant.Response, len(p.DayVariants))
	for i, dv := range p.DayVariants {
		dayVariants[i] = day_variant.ToResponse(dv)
	}
	return FullResponse{
		Response:    toResponse(p),
		DayVariants: dayVariants,
	}
}

func (Plan) TableName() string {
	return "plans"
}
