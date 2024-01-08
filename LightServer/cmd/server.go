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

func (s *apiService) ReportPost(ctx context.Context, req *ReportPostReq) (ReportPostOK, error) {
	s.mux.Lock()
	log.Printf("ReportPost()")
	defer s.mux.Unlock()
}

func (s *apiService) NewError(ctx context.Context, err error) *ErrorStatusCode {
	s.mux.Lock()
	log.Printf("NewError()")
	defer s.mux.Unlock()

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
