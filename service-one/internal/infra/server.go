package infra

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shandler/go-expert-observabilidade/service-one/internal/domain"
	"github.com/shandler/go-expert-observabilidade/service-one/internal/dto"
	"github.com/shandler/go-expert-observabilidade/shared"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Server struct {
	Config  *shared.Config
	ZipCode domain.ZipCode
}

func NewServer(config *shared.Config, zipCode domain.ZipCode) *Server {
	return &Server{
		Config:  config,
		ZipCode: zipCode,
	}
}

func (s *Server) CreateServer() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	// promhttp
	router.Handle("/metrics", promhttp.Handler())

	// request handler
	router.Post("/", s.HandleRequest)

	return router
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := s.Config.OTELTracer.Start(ctx, s.Config.RequestNameOTEL)
	defer span.End()

	request := dto.SearchRequest{}
	json.NewDecoder(r.Body).Decode(&request)

	time.Sleep(20 * time.Millisecond)
	response := s.ZipCode.Search(ctx, request)
	time.Sleep(20 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	json.NewEncoder(w).Encode(response.Body)
}
