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
	"LightServer/arduino"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type apiService struct {
	mux         sync.Mutex
	reports     map[string]api.ReportTuple
	settings    *arduino.Settings
	arduinoPort net.Conn
}

type colorLevels struct {
	threshold int
	red       uint8
	green     uint8
	blue      uint8
}

var palette = [9]colorLevels{
	{threshold: 0, red: 0, green: 0, blue: 0},
	{threshold: 1, red: 0x00, green: 0x8f, blue: 0x00},
	{threshold: 3600 * 1, red: 0x00, green: 0xA0, blue: 0x11},
	{threshold: 3600 * 2, red: 0x8e, green: 0xfa, blue: 0x00},
	{threshold: 3600 * 3, red: 0xff, green: 0xfb, blue: 0x00},
	{threshold: 3600 * 4, red: 0xff, green: 0x93, blue: 0x00},
	{threshold: 3600 * 5, red: 0xff, green: 0x64, blue: 0x00},
	{threshold: 3600 * 6, red: 0xff, green: 0x26, blue: 0x00},
	{threshold: 3600 * 7, red: 0x80, green: 0, blue: 0},
}

func (s *apiService) ResetGet(_ context.Context) (*api.Result, error) {
	return &api.Result{
		Summary: api.OptString{
			Value: "TO BE IMPLEMENTED",
			Set:   true,
		},
	}, nil
}

func (s *apiService) StatusGet(_ context.Context) (*api.Status, error) {
	return &api.Status{}, nil
}

func (s *apiService) ReportPost(_ context.Context, req api.Reports) (*api.Result, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, item := range req {
		if item.Type == api.ReportTupleReportsItem {
			key := item.ReportTuple.Owner + "/" + item.ReportTuple.Repository + "/" + string(item.ReportTuple.Section)
			s.reports[key] = item.ReportTuple
		}
	}

	updateLights(s)

	return &api.Result{
		Summary: api.OptString{
			Value: "Success",
			Set:   true,
		},
	}, nil
}

func updateLights(s *apiService) {
	maxReview := 0
	maxMerge := 0
	maxPull := 0
	for _, item := range s.reports {
		switch item.Section {
		case api.ReportTupleSectionReview:
			{
				if item.Age > maxReview {
					maxReview = item.Age
				}
			}
		case api.ReportTupleSectionMerge:
			{
				if item.Age > maxMerge {
					maxMerge = item.Age
				}
			}
		case api.ReportTupleSectionPull:
			{
				if item.Age > maxPull {
					maxPull = item.Age
				}
			}
		}
	}

	lightCommand := arduino.LightCommand{}
	for i := len(palette) - 1; i >= 0; i-- {
		if maxReview >= palette[i].threshold {
			lightCommand.Red1 = palette[i].red
			lightCommand.Green1 = palette[i].green
			lightCommand.Blue1 = palette[i].blue
			break
		}
	}
	time.Sleep(1 * time.Second)
	for i := len(palette) - 1; i >= 0; i-- {
		if maxMerge >= palette[i].threshold {
			lightCommand.Red2 = palette[i].red
			lightCommand.Green2 = palette[i].green
			lightCommand.Blue2 = palette[i].blue
			break
		}
	}
	time.Sleep(1 * time.Second)
	for i := len(palette) - 1; i >= 0; i-- {
		if maxPull >= palette[i].threshold {
			lightCommand.Red3 = palette[i].red
			lightCommand.Green3 = palette[i].green
			lightCommand.Blue3 = palette[i].blue
			break
		}
	}
	lightCommand.Send(s.arduinoPort)

}

func (s *apiService) NewError(_ context.Context, _ error) *api.ErrorStatusCode {
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

func readSettings() *arduino.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return arduino.ReadSettings(dirname + "/.githubLightBox")
}

func main() {
	service := &apiService{
		reports: map[string]api.ReportTuple{},
	}
	service.settings = readSettings()
	addr := fmt.Sprintf("%s:%d", service.settings.BoxIP, service.settings.BoxPort)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatalf("Could not connect to %s (%v)", addr, err)
		return
	}
	service.arduinoPort = conn

	for i := len(palette) - 1; i >= 0; i-- {
		arduino.LightCommand{
			Red1:   palette[i].red,
			Green1: palette[i].green,
			Blue1:  palette[i].blue,
			Red2:   palette[i].red,
			Green2: palette[i].green,
			Blue2:  palette[i].blue,
			Red3:   palette[i].red,
			Green3: palette[i].green,
			Blue3:  palette[i].blue,
		}.Send(service.arduinoPort)
		time.Sleep(500 * time.Millisecond)
	}

	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatalf("Could not start the server.")
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("Could not listen and serve.")
	}
}
