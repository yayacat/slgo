package handlers

import (
	"bytes"
	"errors"
)

type END struct {
	DataType int
}

// Callback of "END" command, implements handler interface
func (e *END) Callback(client *SeedLinkClient, provider SeedLinkProvider, consumer SeedLinkConsumer, args ...string) error {
	if client.StartTime.IsZero() {
		client.Write([]byte(RES_ERR))
		return errors.New("start time not set")
	}

	var (
		station  = client.Station
		location = client.Location
		network  = client.Network
	)

	// Subscribe to the message queue
	client.Streaming = true
	err := consumer.Subscribe(
		client.RemoteAddr().String(),
		client.Station,
		client.Channels,
		func(data SeedLinkDataPacket) {
			// log.Printf("eventHandler received data: Timestamp=%s, Station=%s, Channel=%s, Data=%v", data.Timestamp, client.Station, data.Channel, data.DataArr)
			newSeq, dataBytes, err := SendSeedLinkPacket(station, location, network, e.DataType, client.GetSequence(), data)
			if err != nil {
				// log.Printf("Error in SendSeedLinkPacket: %v", err)
				consumer.Unsubscribe(client.RemoteAddr().String())
				client.Write([]byte(RES_ERR))
				client.Close()
				return
			}
			// log.Printf("Sending %d bytes to client", len(dataBytes))
			if _, err = client.Write(dataBytes); err != nil {
				consumer.Unsubscribe(client.RemoteAddr().String())
				client.Close()
				// log.Printf("Error writing to client: %v", err)
				return
			}
			client.SetSequence(newSeq)
		},
	)
	if err != nil {
		client.Write([]byte(RES_ERR))
		return err
	}

	// Query history data from database
	historyRecords, err := provider.QueryHistory(client.StartTime, client.EndTime, client.Channels)
	if err != nil {
		client.Write([]byte(RES_ERR))
		return err
	}

	var dataBytesBuf bytes.Buffer
	for _, dataPacket := range historyRecords {
		newSeq, data, err := SendSeedLinkPacket(station, location, network, e.DataType, client.GetSequence(), dataPacket)
		if err != nil {
			client.Write([]byte(RES_ERR))
			return err
		}
		dataBytesBuf.Write(data)
		client.SetSequence(newSeq)
	}
	if _, err = client.Write(dataBytesBuf.Bytes()); err != nil {
		return err
	}

	return nil
}

// Fallback of "END" command, implements handler interface
func (*END) Fallback(client *SeedLinkClient, provider SeedLinkProvider, consumer SeedLinkConsumer, args ...string) {
	client.Close()
}
