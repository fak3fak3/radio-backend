package livego

import (
	"fmt"
	"go-postgres-gorm-gin-api/models"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
)

type LiveGoInstance struct {
	DB *gorm.DB
}

func NewLiveGoLinstance(db *gorm.DB) *LiveGoInstance {
	return &LiveGoInstance{
		DB: db,
	}
}

func (h *LiveGoInstance) Init() {
	ch := make(chan string)
	h.StartPollingStreamStatus(ch)

	go func() {
		for status := range ch {
			fmt.Println(status)
			err := h.DB.Model(&models.StreamData{}).Where("id = ?", 1).Update("status", status).Error
			if err != nil {
				panic(err)
			}
		}
	}()
}

func (h *LiveGoInstance) StartPollingStreamStatus(out chan<- string) {
	go func() {
		for {
			var streamData models.StreamData
			err := h.DB.First(&streamData, 1).Error
			if err != nil {
				panic(err)
			}

			resp, err := http.Get(fmt.Sprintf("http://localhost:7002/live/%s.m3u8", streamData.Room))
			if err != nil {
				log.Println("request failed:", err)
				time.Sleep(1 * time.Second)

				continue
			}
			defer resp.Body.Close()

			spew.Dump(resp.Status)

			if streamData.Status == models.StreamStatusRunning && resp.StatusCode == 403 {
				out <- "empty"

			} else if streamData.Status == models.StreamStatusEmpty && resp.StatusCode == 200 {
				out <- "running"
			}

			time.Sleep(1 * time.Second)
		}
	}()
}
