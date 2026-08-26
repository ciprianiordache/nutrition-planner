package plan

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) CreatePlan(req CreateRequest) (Response, error) { return h.svc.Create(req) }
func (h *Handler) GetPlanFull(id string) (FullResponse, error)    { return h.svc.GetByIDFull(id) }
func (h *Handler) ListPlansByClient(clientID string) ([]Response, error) {
	return h.svc.ListByClientID(clientID)
}
func (h *Handler) UpdatePlan(req UpdateRequest) (Response, error) { return h.svc.Update(req) }
func (h *Handler) DeletePlan(id string) error                     { return h.svc.Delete(id) }
