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

func NewPrayerHandler(service ports.PrayerTimeService) *PrayerHandler {
	return &PrayerHandler{service: service}
}

// GetPrayerTimes godoc
// @Summary Get prayer times and countdown
// @Description Get 5 daily prayer times, current/next status, and countdown to next prayer
// @Tags prayer
// @Produce json
// @Param lat query number true "Latitude"
// @Param long query number true "Longitude"
// @Success 200 {object} domain.PrayerTimesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/prayer-times [get]
func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	latStr := c.Query("lat")
	longStr := c.Query("long")

	if latStr == "" || longStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude is required"})
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
