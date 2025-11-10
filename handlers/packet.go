package handlers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bclswl0827/mseedio"
)

func chunkFloat32Slice(data []float32, chunkSize int) [][]float32 {
	var chunks [][]float32
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func convertFloat32ToInt32(c []float32, multiplier float32) []int32 {
	var int32Data []int32
	for _, v := range c {
		int32Data = append(int32Data, int32(v*multiplier))
	}
	return int32Data
}

func SendSeedLinkPacket(station, location, network string, dataType int, sequence int64, data SeedLinkDataPacket) (newSequence int64, packetBuf []byte, err error) {
	chunks := chunkFloat32Slice(data.DataArr, CHUNK_SIZE)

	dataSpanMs := 1000 / data.SampleRate
	var buf bytes.Buffer

	for i, c := range chunks {
		var miniseed mseedio.MiniSeedData
		miniseed.Init(dataType, mseedio.MSBFIRST)

		startTime := time.UnixMilli(data.Timestamp + int64(i*CHUNK_SIZE*dataSpanMs)).UTC()
		int32Data := convertFloat32ToInt32(c, 1e4)
		// log.Printf("Preparing MiniSEED record: int32Data=%v,NumSamples=%d", int32Data, len(int32Data))
		err := miniseed.Append(int32Data, &mseedio.AppendOptions{
			ChannelCode:    data.Channel,
			StationCode:    station,
			LocationCode:   location,
			NetworkCode:    network,
			SampleRate:     float64(data.SampleRate),
			SequenceNumber: fmt.Sprintf("%06d", sequence),
			StartTime:      startTime,
		})
		if err != nil {
			return 0, nil, err
		}

		// Force 512-byte record
		for i := 0; i < len(miniseed.Series); i++ {
			miniseed.Series[i].BlocketteSection.RecordLength = 9
		}
		slData, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
		if err != nil {
			return 0, nil, err
		}

		slSeq := fmt.Sprintf("SL%06X", sequence)
		buf.Write([]byte(slSeq))
		buf.Write(slData)

		sequence++
	}

	return sequence, buf.Bytes(), nil
}
