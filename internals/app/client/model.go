package client

import "time"

type Client struct {
	ID        string    `db:"id,primary_key,uuid"`
	Name      string    `db:"name,notnull"`
	Weight    float64   `db:"weight,notnull"`
	Height    float64   `db:"height,notnull"`
	Target    string    `db:"target"`
	CreatedAt time.Time `db:"created_at,oncreate"`
	UpdatedAt time.Time `db:"updated_at,onwrite"`
}

// CreateRequest — the data recived from the frontend upon creation
type CreateRequest struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Target string  `json:"target"`
}

// UpdateRequest — the data received from the frontend upon update.
type UpdateRequest struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Target string  `json:"target"`
}

// Response — the data sent to the frontend. Modelul de domeniu (cu tag-uri
// `db`) nu ajunge niciodată direct în UI.
type Response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Weight    float64   `json:"weight"`
	Height    float64   `json:"height"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toResponse(c Client) Response {
	return Response{
		ID:        c.ID,
		Name:      c.Name,
		Weight:    c.Weight,
		Height:    c.Height,
		Target:    c.Target,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (Client) TableName() string {
	return "clients"
}
