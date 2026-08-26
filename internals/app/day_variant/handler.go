package day_variant

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) CreateDayVariant(req CreateRequest) (Response, error) {
	return h.svc.Create(req)
}
func (h *Handler) ListDayVariantsByPlan(planID string) ([]Response, error) {
	return h.svc.ListByPlanIDFull(planID)
}
func (h *Handler) UpdateDayVariant(req UpdateRequest) (Response, error) {
	return h.svc.Update(req)
}
func (h *Handler) DeleteDayVariant(id string) error { return h.svc.Delete(id) }
