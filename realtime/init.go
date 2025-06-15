package realtime

import (
	"encoding/json"
	"fmt"
	"go-postgres-gorm-gin-api/models"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	ch := make(chan string)
	h.StartPollingStreamStatus(ch)

	go func() {
		for status := range ch {
			err := h.DB.Model(&models.StreamData{}).Where("id = ?", 1).Update("status", status).Error
			if err != nil {
				panic(err)
			}
		}
	}()
}

func (h *RealtimeInstance) StartPollingStreamStatus(out chan<- string) {
	go func() {
		for {
			var streamData models.StreamData
			err := h.DB.First(&streamData, 1).Error
			if err != nil {
				panic(err)
			}

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

			var isMainStreamRunning = false

			for _, stream := range data.Streams {
				spew.Dump(stream.Name)
				if stream.Name == "main" {
					isMainStreamRunning = true
				}
			}

			if isMainStreamRunning {
				out <- "running"

			} else {
				out <- "empty"
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
