package realtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type RealtimeInstance struct {
	DB *gorm.DB
}

func NewRealtimeLinstance(db *gorm.DB) *RealtimeInstance {
	return &RealtimeInstance{
		DB: db,
	}
}

func (h *RealtimeInstance) Init() {
	streamsChan := make(chan []Stream)
	h.StartPollingStreamStatus(streamsChan)

	go func() {
		for streams := range streamsChan {
			h.syncStreams(streams)
		}
	}()
}

func (h *RealtimeInstance) StartPollingStreamStatus(out chan<- []Stream) {
	go func() {
		for {
			resp, err := http.Get(fmt.Sprint("http://localhost:1985/api/v1/streams/"))
			if err != nil {
				log.Println("request failed:", err)
				time.Sleep(1 * time.Second)

				continue
			}

			var data StreamResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				log.Println("failed to parse response:", err)
				resp.Body.Close()
				time.Sleep(1 * time.Second)
				continue
			}
			resp.Body.Close()

			out <- data.Streams

			time.Sleep(1 * time.Second)
		}
	}()
}
