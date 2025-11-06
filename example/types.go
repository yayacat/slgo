package main

const (
	SAMPLE_RATE = 50
	TOPIC_NAME  = "seedlink"
)

type adcRawData struct {
	Station    string
	SampleRate int
	Timestamp  int64
	Channel_1  []int32
	Channel_2  []int32
	Channel_3  []int32
}

type eventHandler = func(data *adcRawData)
