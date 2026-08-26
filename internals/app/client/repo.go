package client

import (
	"errors"
	"nutrition-planner/internals/database"

	"github.com/ciprianiordache/crud-depot"
)

type Repo interface {
	Create(c *Client) (string, error)
	GetByID(id string) (*Client, error)
	List(offset, limit int) ([]Client, error)
	Update(c *Client) error
	Delete(id string) error
}

type repo struct {
	depot *crud.CRUD
}

func NewRepo(db *database.Database) Repo { return &repo{depot: crud.New(db, crud.SQLite{})} }

func (r *repo) Create(c *Client) (string, error) {
	return r.depot.Create(c)
}

func (r *repo) GetByID(id string) (*Client, error) {
	var c Client
	err := r.depot.ReadOne(Client{}, "id", id, &c)
	if errors.Is(err, crud.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *repo) List(offset, limit int) ([]Client, error) {
	var clients []Client
	err := r.depot.Get(Client{}, &clients, offset, limit)
	return clients, err
}

func (r *repo) Update(c *Client) error {
	return r.depot.Update(c, "id", c.ID)
}

func (r *repo) Delete(id string) error {
	return r.depot.Delete(&Client{}, "id", id)
}
