package main

import (
	"errors"
	// "log"

	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
	"github.com/yayacat/slgo/handlers"
)

type subscriber struct {
	handler eventHandler
	topic   string
}

type consumer struct {
	messageBus  messagebus.MessageBus
	subscribers cmap.ConcurrentMap[string, subscriber]
}

func (c *consumer) Subscribe(clientId string, station string, channels []handlers.SeedLinkChannel, eventHandler func(handlers.SeedLinkDataPacket), request_id string) error {
	if _, ok := c.subscribers.Get(clientId); ok {
		return errors.New("this client has already subscribed")
	}
	handler := func(data *adcRawData) {
		// Match ADC channels to SeedLink channels
		for _, channel := range channels {
			if channel.ChannelType != "D" {
				continue
			}
			// log.Printf("Matching data: request_id=%s, Station=%s, SeedName=%s with Channel=%s", request_id, data.Station, data.SeedName, channel.ChannelName)
			if data.SeedName == channel.ChannelName {
				go eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  data.Timestamp,
					SampleRate: data.SampleRate,
					Channel:    channel.ChannelName,
					DataArr:    data.Data,
				})
			}
		}
	}
	topic := TOPIC_NAME + "_" + station
	c.subscribers.Set(clientId, subscriber{
		handler: handler,
		topic:   topic,
	})
	c.messageBus.Subscribe(topic, handler)
	return nil
}

func (c *consumer) Unsubscribe(clientId string) error {
	sub, ok := c.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}
	c.messageBus.Unsubscribe(sub.topic, sub.handler)
	c.subscribers.Remove(clientId)
	return nil
}
