package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sctestcase/counter"
)

type Server struct {
	counter counter.Counter
	srv     *http.Server
}

type count struct {
	Count int `json:"count"`
}

const getCount = "/count"

func NewServer(c counter.Counter) (srv *Server) {
	return &Server{
		counter: c,
	}
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == getCount {
		c := s.counter.Count()
		jsonResp := &count{
			Count: c,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResp)
	}

	if r.Method == http.MethodGet {
		s.counter.Inc(r.RequestURI)
	}
}

func (s *Server) AddServer(srv *http.Server) {
	s.srv = srv
}

func (s *Server) Shutdown(ctx context.Context) (err error) {
	var errList []error
	if err := s.srv.Shutdown(ctx); err != nil {
		errList = append(errList, fmt.Errorf("http server shutdown failed: %v", err))
	}

	success := make(chan bool)

	go func(sc chan bool) {
		err = s.counter.Shutdown()
		if err != nil {
			errList = append(errList, fmt.Errorf("save cache: %w", err))
		}

		success <- true
	}(success)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-success:
			if len(errList) > 0 {
				err = fmt.Errorf("%+v", errList)
			}

			return err
		}
	}
}
