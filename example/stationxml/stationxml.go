// Package stationxml 實現 FDSN StationXML 的結構
package stationxml

import (
	"encoding/xml"
	"time"
)

// FDSNStationXML 表示 FDSN StationXML 的根元素
type FDSNStationXML struct {
	XMLName        xml.Name  `xml:"FDSNStationXML"`
	Source         string    `xml:"Source"`
	SchemaVersion  string    `xml:"schemaVersion,attr"`
	Xsi            string    `xml:"xmlns:xsi,attr"`
	Xmlns          string    `xml:"xmlns,attr"`
	SchemaLocation string    `xml:"xsi:schemaLocation,attr"`
	Networks       []Network `xml:"Network"`
}

// Network 表示台網資訊
type Network struct {
	Code        string    `xml:"code,attr"`
	Description string    `xml:"Description"`
	StartDate   string    `xml:"StartDate"`
	Stations    []Station `xml:"Station"`
}

// Station 表示測站資訊
type Station struct {
	Code        string `xml:"code,attr"`
	Description string `xml:"Description"`
	// Latitude     float64   `xml:"Latitude"`
	// Longitude    float64   `xml:"Longitude"`
	Elevation    float64   `xml:"Elevation"`
	Site         Site      `xml:"Site"`
	CreationDate string    `xml:"CreationDate"`
	StartDate    string    `xml:"StartDate"`
	Channels     []Channel `xml:"Channel"`
}

// Site 表示測站地點資訊
type Site struct {
	Name string `xml:"Name"`
}

// Channel 表示頻道資訊
type Channel struct {
	Code     string   `xml:"code,attr"`
	Location string   `xml:"locationCode,attr"`
	Type     []string `xml:"Type"`
	// Latitude     float64  `xml:"Latitude"`
	// Longitude    float64  `xml:"Longitude"`
	Elevation    float64  `xml:"Elevation"`
	Depth        float64  `xml:"Depth"`
	SampleRate   float64  `xml:"SampleRate"`
	Sensor       Sensor   `xml:"Sensor"`
	Response     Response `xml:"Response"`
	CreationDate string   `xml:"CreationDate"`
	StartDate    string   `xml:"StartDate"`
}

// Sensor 表示感應器資訊
type Sensor struct {
	Description string `xml:"Description"`
}

// Response 表示回應資訊
type Response struct {
	InstrumentSensitivity InstrumentSensitivity `xml:"InstrumentSensitivity"`
}

// InstrumentSensitivity 表示儀器靈敏度資訊
type InstrumentSensitivity struct {
	Value       float64 `xml:"Value"`
	Frequency   float64 `xml:"Frequency"`
	InputUnits  Units   `xml:"InputUnits"`
	OutputUnits Units   `xml:"OutputUnits"`
}

// Units 表示單位資訊
type Units struct {
	Name string `xml:"Name"`
}

// DateTime 將時間轉換為 string
func DateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02T15:04:05.000Z")
}
