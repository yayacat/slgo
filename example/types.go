package main

const (
	SAMPLE_RATE = 50
	TOPIC_NAME  = "seedlink"
)

type adcRawData struct {
	Station    string
	SeedName   string
	SampleRate int
	Timestamp  int64
	Data       []float32
}

type eventHandler = func(data *adcRawData)
