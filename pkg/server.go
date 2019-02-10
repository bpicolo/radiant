package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bpicolo/radiant/pkg/backend"
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

func NewServerFromYAML(config string) (*Server, error) {
	store, err := storage.NewYaml(config)
	if err != nil {
		return nil, err
	}
	return &Server{
		manager: backend.NewManager(),
		store:   store,
		engine:  query.NewEngine(),
	}, nil
}

// GetHandler returns the default http handler
func (s *Server) GetHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/query/{queryName}", s.HandleQuery)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/backends", s.GetBackends)

	return r
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

type backendsResponse struct {
	backends []*schema.Backend
}

func (s *Server) GetBackends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write()
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
