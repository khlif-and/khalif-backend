package prayer

import (
	"time"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"
)

type prayerService struct {
	calc *Calculator
}

func NewPrayerService() ports.PrayerTimeService {
	return &prayerService{
		calc: NewCalculator(),
	}
}

// GetPrayerTimes is specific for the Dashboard (Timer + Qibla + Schedule)
func (s *prayerService) GetPrayerTimes(req *domain.PrayerTimesRequest) (*domain.PrayerTimesResponse, error) {
	loc := s.calc.GetTimezone(req.Longitude)
	now := time.Now().In(loc)

	todaySchedule, err := s.calc.GetScheduleForDate(req.Latitude, req.Longitude, now)
	if err != nil {
		return nil, err
	}

	qibla := s.calc.CalculateQibla(req.Latitude, req.Longitude)

	// Helper to format time in HH:mm
	format := func(t time.Time) string {
		return t.In(loc).Format("15:04")
	}

	res := &domain.PrayerTimesResponse{
		Date: now.Format("2006-01-02"),
		Location: domain.Location{
			Lat:  req.Latitude,
			Long: req.Longitude,
		},
		QiblaDirection: qibla,
		Schedule: domain.PrayerSchedule{
			Imsak:   format(s.calc.GetImsak(todaySchedule.Fajr)),
			Subuh:   format(todaySchedule.Fajr),
			Terbit:  format(todaySchedule.Sunrise),
			Dhuha:   format(s.calc.GetDhuha(todaySchedule.Sunrise)),
			Dzuhur:  format(todaySchedule.Zuhr),
			Asar:    format(todaySchedule.Asr),
			Maghrib: format(todaySchedule.Maghrib),
			Isya:    format(todaySchedule.Isha),
		},
	}

	// Calculate Next Prayer and Countdown
	currentSysTime := time.Now().In(loc)

	prayers := []struct {
		Name string
		Time time.Time
	}{
		{"Subuh", todaySchedule.Fajr},
		{"Dzuhur", todaySchedule.Zuhr},
		{"Asar", todaySchedule.Asr},
		{"Maghrib", todaySchedule.Maghrib},
		{"Isya", todaySchedule.Isha},
	}

	var current, next string
	var targetTime time.Time
	foundNext := false

	for i, p := range prayers {
		if p.Time.After(currentSysTime) {
			next = p.Name
			targetTime = p.Time
			
			if i == 0 {
				current = "Isya" // From previous day
			} else {
				current = prayers[i-1].Name
			}
			foundNext = true
			break
		}
	}

	if !foundNext {
		// If past Isya today, next is Subuh tomorrow
		current = "Isya"
		next = "Subuh"
		
		tomorrow := now.Add(24 * time.Hour)
		tomorrowSchedule, err := s.calc.GetScheduleForDate(req.Latitude, req.Longitude, tomorrow)
		
		if err == nil {
			targetTime = tomorrowSchedule.Fajr
		} else {
			targetTime = tomorrow // fallback
		}
	}
	
	diff := targetTime.Sub(currentSysTime)
	if diff < 0 {
		diff = 0
	}
	
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60

	res.Status = domain.PrayerStatus{
		Current: current,
		Next:    next,
		TargetTime: format(targetTime),
		TimeRemaining: domain.TimeRemaining{
			Hours:   hours,
			Minutes: minutes,
			Seconds: seconds,
		},
	}

	return res, nil
}

// GetDailyPrayerTimes is for the Full List View (Clean separate endpoint)
func (s *prayerService) GetDailyPrayerTimes(req *domain.PrayerTimesRequest) (*domain.PrayerTimesResponse, error) {
	// Reuses the logic but optimized for the List View (e.g. maybe no countdown needed? or just consistent?)
	// For now, it returns the same structure but semantically separated.
	// User requested "bedain supaya ga pusing", so having a dedicated method allows future customization (e.g. different response struct)
	// Currently reuse the main Logic to ensure consistency.
	return s.GetPrayerTimes(req)
}
