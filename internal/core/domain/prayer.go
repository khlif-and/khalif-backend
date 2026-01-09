package domain

// PrayerTimesRequest represents the query parameters for prayer times
type PrayerTimesRequest struct {
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"long" binding:"required"`
}

// PrayerSchedule represents the schedule of 5 daily prayers
type PrayerSchedule struct {
	Subuh   string `json:"subuh"`
	Dzuhur  string `json:"dzuhur"`
	Asar    string `json:"asar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

// TimeRemaining represents duplicate countdown data for ease of use
type TimeRemaining struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

// PrayerStatus represents the current state relative to prayer times
type PrayerStatus struct {
	Current       string        `json:"current"`        // e.g. "Asar" (already passed)
	Next          string        `json:"next"`           // e.g. "Maghrib" (upcoming)
	TimeRemaining TimeRemaining `json:"time_remaining"` // Duration until Next
	TargetTime    string        `json:"target_time"`    // Time of Next prayer (e.g. "18:15")
}

// PrayerTimesResponse is the main response object
type PrayerTimesResponse struct {
	Date     string         `json:"date"`
	Location Location       `json:"location"`
	Schedule PrayerSchedule `json:"schedule"`
	Status   PrayerStatus   `json:"status"`
}

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
