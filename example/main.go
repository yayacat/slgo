package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bclswl0827/slgo"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

const (
	HOST = "0.0.0.0"
	PORT = 18000
)

func main() {
	messageBus := messagebus.New(65535)

	// log.Println("test this server with Swarm client: https://volcanoes.usgs.gov/software/swarm/download.shtml")
	log.Printf("starting SeedLink server on %s:%d", HOST, PORT)

	p, err := NewProvider("g:\\work_spaces\\slgo\\example\\stations.json")
	if err != nil {
		log.Fatalf("failed to create provider: %v", err)
	}

	// Create a new SeedLink server with the provider, consumer, and hooks implementations
	server := slgo.New(
		p,
		&consumer{
			messageBus:  messageBus,
			subscribers: cmap.New[subscriber](),
		},
		&hooks{},
	)

	go func() {
		http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
				return
			}

			var data adcRawData
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if data.Station == "" {
				http.Error(w, "station is empty", http.StatusBadRequest)
				return
			}
			topic := TOPIC_NAME + "_" + data.Station
			messageBus.Publish(topic, &data)
			p.mutex.Lock()
			p.hasReceivedRealData = true
			p.mutex.Unlock()
			w.WriteHeader(http.StatusOK)
		})

		http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				stations := p.GetStationsData()
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(stations); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			case http.MethodPost:
				var stations []Station
				if err := json.NewDecoder(r.Body).Decode(&stations); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				p.UpdateStations(stations)
				w.WriteHeader(http.StatusOK)
			default:
				http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			}
		})

		log.Println("starting HTTP server on port 18080")
		if err := http.ListenAndServe(":18080", nil); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	err = server.Start(ctx, HOST, PORT, true)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Stop SeedLink server")
	stop()
}
