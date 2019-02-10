package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bpicolo/radiant/pkg/backend"
	"github.com/bpicolo/radiant/pkg/config"
	"github.com/bpicolo/radiant/pkg/query"
	"github.com/bpicolo/radiant/pkg/storage"
	"github.com/gorilla/mux"
)

// Server is the implementation of the radiant server
type Server struct {
	manager *backend.Manager
	store   storage.Store
	engine  *query.Engine
}

func (s *Server) Shutdown() {
	s.manager.Stop()
}

// Serve starts the server with defaults
func (s *Server) Serve(bind string) {
	srv := &http.Server{
		Addr:         bind,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.GetHandler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	s.Shutdown()
	srv.Shutdown(ctx)

	os.Exit(0)
}

func NewServer(cfg *config.RadiantConfig) (*Server, error) {
	store, err := storage.NewStatic(cfg)
	if err != nil {
		return nil, err
	}
	mgr := backend.NewManager()
	for _, backend := range cfg.Backends {
		err := mgr.AddBackend(backend)
		if err != nil {
			log.Println("Error adding backend", err)
		}
	}

	return &Server{
		manager: mgr,
		store:   store,
		engine:  query.NewEngine(),
	}, nil
}

// GetHandler returns the default http handler
func (s *Server) GetHandler() *mux.Router {
	r := mux.NewRouter()

	// Transparent ES Proxy
	r.PathPrefix("/").HeadersRegexp("Radiant-Proxy-Backend", ".*").HandlerFunc(s.proxy)

	// Radiant api
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/backends", s.getBackends)

	// High-level query layer
	r.HandleFunc("/query/{queryName}", s.HandleQuery)

	return r
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.Header.Get("Radiant-Proxy-Backend"))
	backend := s.manager.GetBackend(name)
	if backend == nil {
		http.Error(
			w,
			fmt.Sprintf("Backend %s not found", name),
			http.StatusServiceUnavailable,
		)
		return
	}

	backend.Backend().Proxy().ServeHTTP(w, r)
}

// type backendsResponse struct {
// 	backends []*schema.Backend
// }

func (s *Server) getBackends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Write()
}

func (s *Server) HandleQuery(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// queryName := vars["queryName"]

	// alias, _ := s.store.GetAlias(queryName)
	// if alias != nil {
	// 	queryName = alias.Name
	// }

	// query, err := s.store.GetQuery(queryName)
	// if err != nil {
	// 	http.Error(
	// 		w,
	// 		fmt.Sprintf("No query with name `%s` was found", queryName),
	// 		http.StatusNotFound,
	// 	)
	// 	return
	// }

}
