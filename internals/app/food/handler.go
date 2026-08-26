package food

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) CreateFood(req CreateRequest) (Response, error) { return h.svc.Create(req) }
func (h *Handler) ListFoodsByMeal(mealID string) ([]Response, error) {
	return h.svc.ListByMealID(mealID)
}
func (h *Handler) UpdateFood(req UpdateRequest) (Response, error) { return h.svc.Update(req) }
func (h *Handler) DeleteFood(id string) error                     { return h.svc.Delete(id) }
