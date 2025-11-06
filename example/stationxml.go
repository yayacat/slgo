// FDSN StationXML Go struct (簡化版)
package main

import "encoding/xml"

// 根據 FDSN StationXML 標準，僅列出常用欄位

type FDSNStationXML struct {
	XMLName        xml.Name  `xml:"FDSNStationXML"`
	Source         string    `xml:"Source"`
	SchemaVersion  string    `xml:"schemaVersion,attr"`
	Xsi            string    `xml:"xmlns:xsi,attr"`
	Xmlns          string    `xml:"xmlns,attr"`
	SchemaLocation string    `xml:"xsi:schemaLocation,attr"`
	Networks       []Network `xml:"Network"`
}

type Network struct {
	Code        string    `xml:"code,attr"`
	Description string    `xml:"Description"`
	StartDate   string    `xml:"StartDate"`
	Stations    []Station `xml:"Station"`
}

type Station struct {
	Code         string    `xml:"code,attr"`
	Description  string    `xml:"Description"`
	Latitude     float64   `xml:"Latitude"`
	Longitude    float64   `xml:"Longitude"`
	Elevation    float64   `xml:"Elevation"`
	Site         Site      `xml:"Site"`
	CreationDate string    `xml:"CreationDate"`
	StartDate    string    `xml:"StartDate"`
	Channels     []Channel `xml:"Channel"`
}

type Site struct {
	Name string `xml:"Name"`
}

type Channel struct {
	Code         string   `xml:"code,attr"`
	Location     string   `xml:"locationCode,attr"`
	Type         []string `xml:"Type"`
	Latitude     float64  `xml:"Latitude"`
	Longitude    float64  `xml:"Longitude"`
	Elevation    float64  `xml:"Elevation"`
	Depth        float64  `xml:"Depth"`
	SampleRate   float64  `xml:"SampleRate"`
	Sensor       Sensor   `xml:"Sensor"`
	Response     Response `xml:"Response"`
	CreationDate string   `xml:"CreationDate"`
	StartDate    string   `xml:"StartDate"`
}

type Sensor struct {
	Description string `xml:"Description"`
}

type Response struct {
	InstrumentSensitivity InstrumentSensitivity `xml:"InstrumentSensitivity"`
}

type InstrumentSensitivity struct {
	Value       float64 `xml:"Value"`
	Frequency   float64 `xml:"Frequency"`
	InputUnits  Units   `xml:"InputUnits"`
	OutputUnits Units   `xml:"OutputUnits"`
}

type Units struct {
	Name string `xml:"Name"`
}
