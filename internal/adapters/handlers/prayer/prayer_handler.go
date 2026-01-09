package prayer

import (
	"net/http"
	"strconv"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type PrayerHandler struct {
	service ports.PrayerTimeService
}

func NewPrayerHandler(s ports.PrayerTimeService) *PrayerHandler {
	return &PrayerHandler{service: s}
}

func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	latStr := c.Query("lat")
	longStr := c.Query("long")

	if latStr == "" || longStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}

	req := &domain.PrayerTimesRequest{
		Latitude:  lat,
		Longitude: long,
	}

	res, err := h.service.GetPrayerTimes(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PrayerHandler) GetDailyPrayerTimes(c *gin.Context) {
	// Reusing same logic for now, but exposed as distinct endpoint
	h.GetPrayerTimes(c)
}
