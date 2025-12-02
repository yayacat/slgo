package main

import (
	"context"
	"encoding/json"
	"encoding/xml" // Added
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal" // Added
	"syscall"

	stationxml "github.com/yayacat/slgo/example/stationxml"

	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus" // Added
	"github.com/yayacat/slgo"
)

func fdsnStationHandler(p *provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		fdsnXML := stationxml.FDSNStationXML{
			Source:         p.GetOrganization(),
			SchemaVersion:  "1.0",
			Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
			Xmlns:          "http://www.fdsn.org/xml/station/1",
			SchemaLocation: "http://www.fdsn.org/xml/station/1 http://www.fdsn.org/xml/station/fdsn-station-1.0.xsd",
		}

		for _, s := range p.GetStationsData() {
			network := stationxml.Network{
				Code:        s.Network,
				Description: s.Description,
				StartDate:   stationxml.DateTime(p.GetStartTime()),
			}

			sta := stationxml.Station{
				Code:        s.Station,
				Description: s.Description,
				// Latitude:    s.Latitude,
				// Longitude:   s.Longitude,
				Elevation: 0.0,
				Site: stationxml.Site{
					Name: s.Description,
				},
				CreationDate: stationxml.DateTime(p.GetStartTime()),
				StartDate:    stationxml.DateTime(p.GetStartTime()),
			}

			for _, stream := range s.Streams {
				channel := stationxml.Channel{
					Code:     stream.SeedName,
					Location: stream.Location,
					Type:     []string{stream.Type},
					// Latitude:   s.Latitude,
					// Longitude:  s.Longitude,
					Elevation:  0.0,
					Depth:      0.0,
					SampleRate: SAMPLE_RATE,
					Sensor: stationxml.Sensor{
						Description: s.Description + " " + stream.SeedName,
					},
					Response: stationxml.Response{
						InstrumentSensitivity: stationxml.InstrumentSensitivity{
							Value:       1.0,
							Frequency:   1.0,
							InputUnits:  stationxml.Units{Name: "M/S"},
							OutputUnits: stationxml.Units{Name: "COUNTS"},
						},
					},
					CreationDate: stationxml.DateTime(p.GetStartTime()),
					StartDate:    stationxml.DateTime(p.GetStartTime()),
				}
				sta.Channels = append(sta.Channels, channel)
			}
			network.Stations = append(network.Stations, sta)
			fdsnXML.Networks = append(fdsnXML.Networks, network)
		}

		output, err := xml.MarshalIndent(fdsnXML, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(xml.Header))
		w.Write(output)
	}
}

func main() {
	stationsPath := flag.String("stations", "stations.json", "path to stations.json file")
	flag.Parse()

	messageBus := messagebus.New(99999)

	// log.Println("test this server with Swarm client: https://volcanoes.usgs.gov/software/swarm/download.shtml")
	log.Printf("starting SeedLink server on %s:%d", HOST, PORT)

	p, err := NewProvider(*stationsPath)
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

			// log.Printf("Received data: Station=%s, SeedName=%s, Data=%v", data.Station, data.SeedName, data.Data)

			if data.Station == "" {
				http.Error(w, "station is empty", http.StatusBadRequest)
				return
			}

			if data.SeedName == "" {
				http.Error(w, "seed name is empty", http.StatusBadRequest)
				return
			}

			// log.Printf("Received data: Station=%s, SeedName=%s", data.Station, data.SeedName)
			topic := TOPIC_NAME + "_" + data.Station
			messageBus.Publish(topic, &data)
			p.mutex.Lock()
			p.hasReceivedRealData = true
			p.mutex.Unlock()
			w.WriteHeader(http.StatusOK)
		})

		// http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		// 	switch r.Method {
		// 	case http.MethodGet:
		// 		stations := p.GetStationsData()
		// 		w.Header().Set("Content-Type", "application/json")
		// 		if err := json.NewEncoder(w).Encode(stations); err != nil {
		// 			http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		}
		// 	case http.MethodPost:
		// 		var stations []Station
		// 		if err := json.NewDecoder(r.Body).Decode(&stations); err != nil {
		// 			http.Error(w, err.Error(), http.StatusBadRequest)
		// 			return
		// 		}

		// 		p.UpdateStations(stations)
		// 		w.WriteHeader(http.StatusOK)
		// 	default:
		// 		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		// 	}
		// })

		http.HandleFunc("/fdsnws/station/1/query", fdsnStationHandler(p))

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
