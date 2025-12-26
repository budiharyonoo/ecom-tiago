package orders

import (
	"github/budiharyonoo/ecom-tiago/internal/utils/request"
	"github/budiharyonoo/ecom-tiago/internal/utils/response"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(s Service) *handler {
	return &handler{
		service:   s,
		validator: validator.New(),
	}
}

func (h *handler) Store(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := request.Json(r, &req); err != nil {
		response.Json(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Json(w, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	h.service.Store(r.Context(), req)

	response.Json(w, http.StatusCreated)
}
