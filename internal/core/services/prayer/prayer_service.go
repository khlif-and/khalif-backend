package prayer

import (
	"errors"
	"math"
	"time"

	"khalif-backend/internal/core/domain"
	"khalif-backend/internal/core/ports"

	"github.com/hablullah/go-prayer"
)

type prayerService struct{}

func NewPrayerService() ports.PrayerTimeService {
	return &prayerService{}
}

func (s *prayerService) GetPrayerTimes(req *domain.PrayerTimesRequest) (*domain.PrayerTimesResponse, error) {
	// Calculate timezone offest based on longitude
	zoneOffset := int(math.Round(req.Longitude / 15.0))
	loc := time.FixedZone("Local", zoneOffset*3600)
	now := time.Now().In(loc)
	year := now.Year()

	// Calculate prayer times for the whole year
	schedules, err := prayer.Calculate(
		prayer.Config{
			Latitude:           req.Latitude,
			Longitude:          req.Longitude,
			Timezone:           loc,
			TwilightConvention: prayer.Kemenag(),
			AsrConvention:      prayer.Shafii,
			PreciseToSeconds:   true,
		},
		year,
	)
	if err != nil {
		return nil, err
	}

	// Find today's schedule
	var todaySchedule prayer.Schedule
	found := false
	targetDate := now.Format("2006-01-02")

	for _, s := range schedules {
		// s.Date is a string in format "2006-01-02"
		if s.Date == targetDate {
			todaySchedule = s
			found = true
			break
		}
	}

	if !found {
		// Fallback for timezone edge cases (date shift)
		// Or errors.
		return nil, errors.New("schedule not found for today")
	}

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
		Schedule: domain.PrayerSchedule{
			Subuh:   format(todaySchedule.Fajr),
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
		
		var tomorrowSchedule prayer.Schedule
		foundTomorrow := false
		targetTom := tomorrow.Format("2006-01-02")
		
		if tomorrow.Year() == year {
			for _, s := range schedules {
				if s.Date == targetTom {
					tomorrowSchedule = s
					foundTomorrow = true
					break
				}
			}
		} else {
			// Calculate for next year
			nextYearSchedules, _ := prayer.Calculate(
				prayer.Config{
					Latitude:           req.Latitude,
					Longitude:          req.Longitude,
					Timezone:           loc,
					TwilightConvention: prayer.Kemenag(),
					AsrConvention:      prayer.Shafii,
					PreciseToSeconds:   true,
				},
				tomorrow.Year(),
			)
			for _, s := range nextYearSchedules {
				if s.Date == targetTom {
					tomorrowSchedule = s
					foundTomorrow = true
					break
				}
			}
		}

		if foundTomorrow {
			targetTime = tomorrowSchedule.Fajr
		} else {
			// Fallback
			targetTime = tomorrow 
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
