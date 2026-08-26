package client

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateClient(req CreateRequest) (Response, error) {
	return h.svc.Create(req)
}

func (h *Handler) GetClient(id string) (Response, error) {
	return h.svc.GetByID(id)
}

func (h *Handler) ListClients(offset, limit int) ([]Response, error) {
	return h.svc.List(offset, limit)
}

func (h *Handler) UpdateClient(req UpdateRequest) (Response, error) {
	return h.svc.Update(req)
}

func (h *Handler) Delete(id string) error {
	return h.svc.Delete(id)
}
