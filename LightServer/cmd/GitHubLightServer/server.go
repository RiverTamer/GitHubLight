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
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type apiService struct {
	mux           sync.Mutex
	clientReports map[string]api.ClientReport
	settings      *arduino.Settings
	arduinoPort   net.Conn
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
	clear(s.clientReports)
	updateLights(s)

	return &api.Result{
		Summary: api.OptString{
			Value: "TO BE IMPLEMENTED",
			Set:   true,
		},
	}, nil
}

func (s *apiService) StatusGet(_ context.Context) (*api.Status, error) {
	status := api.Status{
		Reports: make([]api.ReportsItem, 0),
	}
	for _, report := range s.clientReports {
		status.Reports = append(status.Reports, report.Reports...)
	}
	return &status, nil
}

func (s *apiService) ReportPost(_ context.Context, req *api.ClientReport) (*api.Result, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.clientReports[req.Clientid] = *req

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
	pullCount := 0
	for _, clientReport := range s.clientReports {
		for _, item := range clientReport.Reports {
			if item.Type == api.ReportTupleReportsItem {
				switch item.ReportTuple.Section {
				case api.ReportTupleSectionReview:
					{
						if item.ReportTuple.Age > maxReview {
							maxReview = item.ReportTuple.Age
						}
					}
				case api.ReportTupleSectionMerge:
					{
						if item.ReportTuple.Age > maxMerge {
							maxMerge = item.ReportTuple.Age
						}
					}
				case api.ReportTupleSectionPull:
					{
						if item.ReportTuple.Age > 0 {
							pullCount = pullCount + item.ReportTuple.Age
						}
					}
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
	if maxMerge == 1 {
		lightCommand.Blue2 = 40
	} else {
		for i := len(palette) - 1; i >= 0; i-- {
			if maxMerge >= palette[i].threshold {
				lightCommand.Red2 = palette[i].red
				lightCommand.Green2 = palette[i].green
				lightCommand.Blue2 = palette[i].blue
				break
			}
		}
	}
	if pullCount < 25 {
		lightCommand.Blue3 = uint8(pullCount * 10)
	} else if pullCount < 50 {
		lightCommand.Green3 = 0xff
	} else if pullCount < 75 {
		lightCommand.Red3 = 0x80
		lightCommand.Green3 = 0x80
	} else {
		lightCommand.Red3 = 0xff
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
		clientReports: map[string]api.ClientReport{},
	}
	service.settings = readSettings()
	addr := fmt.Sprintf("%s:%d", service.settings.BoxIP, service.settings.BoxPort)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatalf("Could not connect to %s (%v)", addr, err)
		return
	}
	service.arduinoPort = conn

	for i := len(palette); i >= 0; i-- {
		g1 := rand.Intn(len(palette))
		g2 := rand.Intn(len(palette))
		g3 := rand.Intn(len(palette))
		arduino.LightCommand{
			Red1:   palette[g1].red,
			Green1: palette[g1].green,
			Blue1:  palette[g1].blue,
			Red2:   palette[g2].red,
			Green2: palette[g2].green,
			Blue2:  palette[g2].blue,
			Red3:   palette[g3].red,
			Green3: palette[g3].green,
			Blue3:  palette[g3].blue,
		}.Send(service.arduinoPort)
		time.Sleep(500 * time.Millisecond)
	}
	arduino.LightCommand{
		Red1:   0,
		Green1: 0,
		Blue1:  0,
		Red2:   0,
		Green2: 0,
		Blue2:  0,
		Red3:   0,
		Green3: 0,
		Blue3:  0,
	}.Send(service.arduinoPort)

	srv, err := api.NewServer(service)
	if err != nil {
		log.Fatalf("Could not start the server.")
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("Could not listen and serve.")
	}
}
