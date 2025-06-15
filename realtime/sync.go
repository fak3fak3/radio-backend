package realtime

import (
	"go-postgres-gorm-gin-api/models"
	"time"
)

func (h *RealtimeInstance) syncStreams(incoming []Stream) error {
	var existing []models.StreamData
	if err := h.DB.Find(&existing).Error; err != nil {
		return err
	}

	existingMap := make(map[string]models.StreamData)
	for _, s := range existing {
		existingMap[s.Room] = s
	}

	incomingMap := make(map[string]Stream)
	for _, s := range incoming {
		incomingMap[s.Name] = s
	}

	// добавление и обновление
	for _, s := range incoming {
		status := models.StreamStatusCreated

		existing, found := existingMap[s.Name]
		if !found || existing.Status != status || existing.Kbps != s.Kbps.Send30s || existing.Latency != s.LiveMS {

			record := models.StreamData{
				Room:    s.Name,
				Status:  status,
				Kbps:    s.Kbps.Send30s,
				Latency: s.LiveMS,
			}
			if found {
				record.ID = existing.ID
			}
			if err := h.DB.Save(&record).Error; err != nil {
				return err
			}
		}
	}

	for room, s := range existingMap {
		if _, ok := incomingMap[room]; !ok {
			if s.Status != models.StreamStatusEnded {
				if err := h.DB.Model(&models.StreamData{}).
					Where("id = ?", s.ID).
					Updates(map[string]interface{}{
						"status":     models.StreamStatusEnded,
						"updated_at": time.Now(),
					}).Error; err != nil {
					return err
				}
			}
			if err := h.DB.Delete(&models.StreamData{}, s.ID).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
