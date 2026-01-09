package domain

type PrayerTimesRequest struct {
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"long" binding:"required"`
}

type PrayerSchedule struct {
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Terbit  string `json:"terbit"` // Sunrise
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Asar    string `json:"asar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

type TimeRemaining struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

type PrayerStatus struct {
	Current       string        `json:"current"`
	Next          string        `json:"next"`
	TargetTime    string        `json:"target_time"` // HH:mm of the next prayer
	TimeRemaining TimeRemaining `json:"time_remaining"`
}

type PrayerTimesResponse struct {
	Date           string         `json:"date"`
	Location       Location       `json:"location"`
	QiblaDirection float64        `json:"qibla_direction"` // Degrees from North
	Schedule       PrayerSchedule `json:"schedule"`
	Status         PrayerStatus   `json:"status,omitempty"` // Optional for full list view
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
