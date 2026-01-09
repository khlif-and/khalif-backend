package prayer

import (
	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
	"khalif-backend/pkg/messages"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrayerHandler struct {
	service ports.PrayerTimeService
}

func NewPrayerHandler(service ports.PrayerTimeService) *PrayerHandler {
	return &PrayerHandler{service: service}
}

func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	latStr := c.Query("lat")
	longStr := c.Query("long")

	if latStr == "" || longStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrLatLongRequired})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidLatitude})
		return
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidLongitude})
		return
	}

	req := &domain.PrayerTimesRequest{
		Latitude:  lat,
		Longitude: long,
	}

	res, err := h.service.GetPrayerTimes(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrInternalServer})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PrayerHandler) GetDailyPrayerTimes(c *gin.Context) {
	h.GetPrayerTimes(c)
}
