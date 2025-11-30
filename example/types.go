package main

const (
	SAMPLE_RATE = 50
	TOPIC_NAME  = "seedlink"
	HOST        = "0.0.0.0"
	PORT        = 18000
)

type adcRawData struct {
	Station    string
	SeedName   string
	SampleRate int
	Timestamp  int64
	Data       []float32
}

type eventHandler = func(data *adcRawData)
