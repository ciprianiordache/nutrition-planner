package supplement

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateSupplement(req CreateRequest) (Response, error) {
	return h.svc.Create(req)
}

func (h *Handler) ListSupplementsByMeal(mealID string) ([]Response, error) {
	return h.svc.ListByMealID(mealID)
}

func (h *Handler) UpdateSupplement(req UpdateRequest) (Response, error) {
	return h.svc.Update(req)
}

func (h *Handler) DeleteSupplement(id string) error {
	return h.svc.Delete(id)
}
