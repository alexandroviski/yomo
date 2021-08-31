package main

import (
	"io"
	"log"
	"math/rand"
	"time"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo"
)

type noiseData struct {
	Noise float32 `y3:"0x11"` // Noise value
	Time  int64   `y3:"0x12"` // Timestamp (ms)
	From  string  `y3:"0x13"` // Source IP
}

func main() {
	// connect to YoMo-Zipper.
	source := yomo.NewSource(yomo.WithName("yomo-source"), yomo.WithZipperEndpoint("localhost:9000"))
	err := source.Connect()
	if err != nil {
		log.Printf("[source] ❌ Emit the data to YoMo-Zipper failure with err: %v", err)
		return
	}

	defer source.Close()

	source.SetDataTag(0x33)
	// generate mock data and send it to YoMo-Zipper in every 100 ms.
	generateAndSendData(source)
}

var codec = y3.NewCodec(0x10)

func generateAndSendData(stream io.Writer) {
	for {
		// generate random data.
		data := noiseData{
			Noise: rand.New(rand.NewSource(time.Now().UnixNano())).Float32() * 200,
			Time:  time.Now().UnixNano() / int64(time.Millisecond),
			From:  "localhost",
		}

		// Encode data via Y3 codec https://github.com/yomorun/y3-codec.
		sendingBuf, _ := codec.Marshal(data)

		// send data via QUIC stream.
		_, err := stream.Write(sendingBuf)
		if err != nil {
			log.Printf("[source] ❌ Emit %v to YoMo-Zipper failure with err: %v", data, err)
		} else {
			log.Printf("[source] ✅ Emit %v to YoMo-Zipper", data)
		}

		time.Sleep(100 * time.Millisecond)
	}
}