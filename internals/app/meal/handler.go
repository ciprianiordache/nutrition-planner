package meal

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) CreateMeal(req CreateRequest) (Response, error) { return h.svc.Create(req) }
func (h *Handler) ListMealsByDayVariant(dayVariantID string) ([]Response, error) {
	return h.svc.ListByDayVariantIDFull(dayVariantID)
}
func (h *Handler) UpdateMeal(req UpdateRequest) (Response, error) { return h.svc.Update(req) }
func (h *Handler) DeleteMeal(id string) error                     { return h.svc.Delete(id) }
