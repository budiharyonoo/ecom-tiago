package products

import (
	"github/budiharyonoo/ecom-tiago/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.Context())
	if err != nil {
		response.Json(w, http.StatusInternalServerError)
		return
	}

	response.Json(w, http.StatusOK, products)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	productId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Json(w, http.StatusBadRequest)
		return
	}

	product, err := h.service.GetById(r.Context(), uint64(productId))
	if err != nil {
		response.Json(w, http.StatusNotFound)
		return
	}

	response.Json(w, http.StatusOK, product)
}
