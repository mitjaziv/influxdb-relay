package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	client := influxdb2.NewClientWithOptions("http://localhost:9096", "my-token",
		influxdb2.DefaultOptions().SetBatchSize(20))

	// Get non-blocking write client.
	writeAPI := client.WriteAPI("my-org", "my-bucket")

	// Create and write points.
	for i := 0; i < 100; i++ {
		p := influxdb2.NewPoint(
			"system",
			map[string]string{
				"id":       fmt.Sprintf("rack_%v", i%10),
				"vendor":   "AWS",
				"hostname": fmt.Sprintf("host_%v", i%100),
			},
			map[string]interface{}{
				"temperature": rand.Float64() * 80.0,
				"disk_free":   rand.Float64() * 1000.0,
				"disk_total":  (i/10 + 1) * 1000000,
				"mem_total":   (i/100 + 1) * 10000000,
				"mem_free":    rand.Uint64(),
			},
			time.Now())

		// Write points asynchronously.
		writeAPI.WritePoint(p)
	}

	// Force all unwritten data to be sent.
	writeAPI.Flush()

	// Ensures background processes finishes.
	client.Close()
}
