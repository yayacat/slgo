package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/bclswl0827/slgo/handlers"
)

type Stream struct {
	SeedName string `json:"seedname"`
	Location string `json:"location"`
	Type     string `json:"type"`
}

type Station struct {
	Station     string   `json:"station"`
	Network     string   `json:"network"`
	Description string   `json:"description"`
	Streams     []Stream `json:"streams"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
}

type provider struct {
	startTime         time.Time
	stations          []Station
	mutex             sync.RWMutex
	hasReceivedRealData bool
}

func NewProvider(filePath string) (*provider, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var stations []Station
	err = json.Unmarshal(file, &stations)
	if err != nil {
		return nil, err
	}

	return &provider{
		startTime: time.Now().UTC(),
		stations:  stations,
		mutex:     sync.RWMutex{},
		hasReceivedRealData: false,
	}, nil
}

func (p *provider) UpdateStations(stations []Station) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.stations = stations
}

func (p *provider) GetStationsData() []Station {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.stations
}

func (p *provider) GetSoftware() string {
	return "slgo"
}

func (p *provider) GetStartTime() time.Time {
	return p.startTime
}

func (p *provider) GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func (p *provider) GetOrganization() string {
	return "anyshake.org"
}

func (p *provider) GetStations() []handlers.SeedLinkStation {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var slStations []handlers.SeedLinkStation
	for _, s := range p.stations {
		slStations = append(slStations, handlers.SeedLinkStation{
			BeginSequence: "000000",
			EndSequence:   "FFFFFF",
			Station:       s.Station,
			Network:       s.Network,
			Description:   s.Description,
			Latitude:      s.Latitude,
			Longitude:     s.Longitude,
		})
	}
	return slStations
}

func (p *provider) GetStreams() []handlers.SeedLinkStream {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var slStreams []handlers.SeedLinkStream
	for _, s := range p.stations {
		for _, st := range s.Streams {
			slStreams = append(slStreams, handlers.SeedLinkStream{
				BeginTime: p.GetStartTime().Format("2006-01-02 15:04:01"),
				EndTime:   p.GetCurrentTime().Format("2006-01-02 15:04:01"),
				SeedName:  st.SeedName,
				Location:  st.Location,
				Type:      st.Type,
				Station:   s.Station,
			})
		}
	}
	return slStreams
}

func (p *provider) GetCapabilities() []handlers.SeedLinkCapability {
	return []handlers.SeedLinkCapability{
		{Name: "info:all"}, {Name: "info:gaps"}, {Name: "info:streams"},
		{Name: "dialup"}, {Name: "info:id"}, {Name: "multistation"},
		{Name: "window-extraction"}, {Name: "info:connections"},
		{Name: "info:capabilities"}, {Name: "info:stations"},
	}
}

func (p *provider) QueryHistory(startTime, endTime time.Time, channels []handlers.SeedLinkChannel) ([]handlers.SeedLinkDataPacket, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.hasReceivedRealData {
		return []handlers.SeedLinkDataPacket{}, nil
	}

	var dataPackets []handlers.SeedLinkDataPacket

	// Generate random data packets for each channel, every second
	startTimestamp, endTimestamp := startTime.UnixMilli(), endTime.UnixMilli()
	for i := startTimestamp; i < endTimestamp; i += 1000 {
		for _, channel := range channels {
			dataPacket := handlers.SeedLinkDataPacket{
				Timestamp:  i,
				SampleRate: SAMPLE_RATE,
				Channel:    channel.ChannelName,
				DataArr:    generateRandomArray(SAMPLE_RATE, -32768, 32768),
			}
			dataPackets = append(dataPackets, dataPacket)
		}
	}

	return dataPackets, nil
}
