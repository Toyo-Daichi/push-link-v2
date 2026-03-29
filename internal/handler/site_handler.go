package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"push-link-v2/internal/domain"
	"push-link-v2/internal/usecase"
	"push-link-v2/internal/view"
)

type SiteHandler struct {
	renderer    *view.TemplateRenderer
	siteUsecase *usecase.SiteUsecase
}

type IndexPageData struct {
	SiteList domain.SiteListResult
	Filter   domain.SiteFilter
}

func NewSiteHandler(renderer *view.TemplateRenderer, siteUsecase *usecase.SiteUsecase) *SiteHandler {
	return &SiteHandler{
		renderer:    renderer,
		siteUsecase: siteUsecase,
	}
}

func (h *SiteHandler) Index(w http.ResponseWriter, r *http.Request) {
	filter := requestFilter(r)
	result, err := h.siteUsecase.List(filter)
	if err != nil {
		http.Error(w, "failed to load sites", http.StatusInternalServerError)
		return
	}

	data := IndexPageData{
		SiteList: result,
		Filter:   normalizeFilter(filter, result.Strategy),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "layout", data); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func (h *SiteHandler) SiteListPartial(w http.ResponseWriter, r *http.Request) {
	filter := requestFilter(r)
	result, err := h.siteUsecase.List(filter)
	if err != nil {
		http.Error(w, "failed to load sites", http.StatusInternalServerError)
		return
	}

	data := IndexPageData{
		SiteList: result,
		Filter:   normalizeFilter(filter, result.Strategy),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.renderer.Render(w, "site-list", data); err != nil {
		http.Error(w, "failed to render partial", http.StatusInternalServerError)
	}
}

func (h *SiteHandler) SiteListAPI(w http.ResponseWriter, r *http.Request) {
	filter := requestFilter(r)
	result, err := h.siteUsecase.List(filter)
	if err != nil {
		http.Error(w, "failed to load sites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
	}
}

func (h *SiteHandler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func requestFilter(r *http.Request) domain.SiteFilter {
	return domain.SiteFilter{
		Strategy: strings.TrimSpace(r.URL.Query().Get("strategy")),
		Query:    strings.TrimSpace(r.URL.Query().Get("q")),
		Tag:      strings.TrimSpace(r.URL.Query().Get("tag")),
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
	}
}

func normalizeFilter(filter domain.SiteFilter, strategy string) domain.SiteFilter {
	filter.Strategy = strategy
	return filter
}
