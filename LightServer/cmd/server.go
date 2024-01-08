//
//  server.go
//  LightServer
//
//  Created by Karl Kraft on 1/7/2024
//  Copyright 2024 Karl Kraft. All rights reserved
//

package main

import (
	"LightServer/api"
	"context"
	"log"
	"net/http"
	"sync"
)

type apiService struct {
	mux sync.Mutex
}

func (s *apiService) ReportPost(ctx context.Context, req api.Reports) (*api.Result, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return &api.Result{
		Summary: api.OptString{
			Value: "Success",
			Set:   true,
		},
	}, nil
}

func (s *apiService) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	s.mux.Lock()
	log.Printf("NewError()")
	defer s.mux.Unlock()
	return &api.ErrorStatusCode{
		StatusCode: 404,
		Response: api.Error{
			Summary: api.OptString{
				Value: "Generic error",
				Set:   true,
			},
		},
	}

}

func main() {
	service := &apiService{}
	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatalf("Could not start the server.")
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("Could not listen and serve.")
	}
}
