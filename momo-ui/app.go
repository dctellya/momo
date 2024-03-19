package main

import (
	"changeme/components"
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func NewChiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/initial", templ.Handler(components.Pages([]struct {
		Path  string
		Label string
	}{
		{"/greet", "Greet form"},
		{"/modal", "Modal"},
		{"/graphcontent", "Graph"},
		{"/gaugecontent", "Gauge Graph"},
	}, struct {
		Version string
		Text    string
	}{
		version, "No update available",
	})).ServeHTTP)
	r.Get("/greet", templ.Handler(components.GreetForm("/greet")).ServeHTTP)
	r.Get("/graphcontent", func(w http.ResponseWriter, r *http.Request) {
		w.Write(loadChart())
	})
	r.Get("/gaugecontent", func(w http.ResponseWriter, r *http.Request) {
		w.Write(loadGaugeChart())
	})
	r.Post("/greet", components.Greet)
	r.Get("/modal", templ.Handler(components.TestPage("#modal", "outerHTML")).ServeHTTP)
	r.Post("/modal", templ.Handler(components.ModalPreview("Title for the modal", "Sample Data")).ServeHTTP)
	return r
}

func loadChart() []byte {
	g := GraphExamples{}
	return g.LoadChart()
}

func loadGaugeChart() []byte {
	g := GraphExamples{}
	return g.LoadGaugeChart()
}
