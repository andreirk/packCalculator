package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"packCalculator/server/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PackHandler struct {
	packService *service.PackService
}

func NewPackHandler(packService *service.PackService) *PackHandler {
	return &PackHandler{
		packService: packService,
	}
}

type CalculatePacksRequest struct {
	ItemsOrdered int `json:"items_ordered"`
}

type CalculatePacksResponse struct {
	Packs map[int]int `json:"packs"`
}

type UpdatePackSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

func (h *PackHandler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	itemsStr := r.URL.Query().Get("items")
	items, err := strconv.Atoi(itemsStr)
	if err != nil {
		http.Error(w, "Invalid items parameter", http.StatusBadRequest)
		return
	}

	fmt.Println("Get items here", items)
	packs := h.packService.CalculatePacks(items)
	fmt.Println("Get packs here after", packs)

	response := CalculatePacksResponse{
		Packs: packs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PackHandler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes := h.packService.GetCurrentSizes()
	respondJSON(w, http.StatusOK, map[string]interface{}{"sizes": sizes})
}

func (h *PackHandler) AddPackSize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Size int `json:"size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Size <= 0 {
		respondError(w, http.StatusBadRequest, "Pack size must be positive")
		return
	}

	if err := h.packService.AddPackSize(req.Size); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to add pack size")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"message": "Pack size added"})
}

func (h *PackHandler) RemovePackSize(w http.ResponseWriter, r *http.Request) {
	sizeStr := chi.URLParam(r, "size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid size parameter")
		return
	}

	if err := h.packService.RemovePackSize(size); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to remove pack size")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"message": "Pack size removed"})
}

func (h *PackHandler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Sizes []int `json:"sizes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Sizes) == 0 {
		respondError(w, http.StatusBadRequest, "At least one pack size is required")
		return
	}

	for _, size := range req.Sizes {
		if size <= 0 {
			respondError(w, http.StatusBadRequest, "All pack sizes must be positive")
			return
		}
	}

	if err := h.packService.UpdatePackSizes(req.Sizes); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update pack sizes")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"message": "Pack sizes updated"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
