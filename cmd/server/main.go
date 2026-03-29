package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"push-link-v2/internal/handler"
	"push-link-v2/internal/repository"
	"push-link-v2/internal/service"
	"push-link-v2/internal/usecase"
	"push-link-v2/internal/view"
)

func main() {
	templateRenderer, err := view.NewTemplateRenderer("web/templates/*.html")
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	siteRepository := repository.NewMemorySiteRepository(repository.SeedSites())
	strategyRegistry := service.NewSiteStrategyRegistry()
	siteUsecase := usecase.NewSiteUsecase(siteRepository, strategyRegistry)
	siteHandler := handler.NewSiteHandler(templateRenderer, siteUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/", siteHandler.Index)
	mux.HandleFunc("/healthz", siteHandler.Healthz)
	mux.HandleFunc("/ui/sites", siteHandler.SiteListPartial)
	mux.HandleFunc("/api/v1/sites", siteHandler.SiteListAPI)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	server := &http.Server{
		Addr:              ":" + port(),
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server started on http://localhost:%s", port())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}

func port() string {
	if value := os.Getenv("PORT"); value != "" {
		return value
	}
	return "8080"
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
