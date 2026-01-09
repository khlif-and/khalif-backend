package prayer

import (
	"errors"
	"math"
	"time"

	"github.com/hablullah/go-prayer"
)

// Calculator handles math and library interactions
type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateQibla returns the Qibla direction in degrees from North
// Kaaba coordinates: 21.422487, 39.826206
func (c *Calculator) CalculateQibla(lat, long float64) float64 {
	kaabaLat := 21.422487
	kaabaLong := 39.826206

	// Convert to radians
	latRad := lat * math.Pi / 180
	longRad := long * math.Pi / 180
	kaabaLatRad := kaabaLat * math.Pi / 180
	kaabaLongRad := kaabaLong * math.Pi / 180

	y := math.Sin(kaabaLongRad-longRad) * math.Cos(kaabaLatRad)
	x := math.Cos(latRad)*math.Sin(kaabaLatRad) - math.Sin(latRad)*math.Cos(kaabaLatRad)*math.Cos(kaabaLongRad-longRad)

	bearingRad := math.Atan2(y, x)
	bearingDeg := bearingRad * 180 / math.Pi

	// Normalize to 0-360
	return math.Mod(bearingDeg+360, 360)
}

func (c *Calculator) GetScheduleForDate(lat, long float64, date time.Time) (*prayer.Schedule, error) {
	// Calculate timezone offest based on longitude
	zoneOffset := int(math.Round(long / 15.0))
	loc := time.FixedZone("Local", zoneOffset*3600)
	
	year := date.Year()

	// library calculates for the whole year
	schedules, err := prayer.Calculate(
		prayer.Config{
			Latitude:           lat,
			Longitude:          long,
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

	targetDate := date.Format("2006-01-02")
	for i := range schedules {
		if schedules[i].Date == targetDate {
			return &schedules[i], nil
		}
	}

	return nil, errors.New("schedule not found for date")
}

// Helper methods to calculate derived times
func (c *Calculator) GetImsak(fajr time.Time) time.Time {
	// Kemenag standard: Imsak is 10 minutes before Shubuh
	return fajr.Add(-10 * time.Minute)
}

func (c *Calculator) GetDhuha(sunrise time.Time) time.Time {
	// Kemenag/Standard: Dhuha starts when sun is approx 7 degrees up, 
	// typically ~20-25 mins after sunrise. We'll use +20 mins as a safe approximation for "Syuruq to Dhuha"
	return sunrise.Add(20 * time.Minute)
}

func (c *Calculator) GetTimezone(long float64) *time.Location {
	zoneOffset := int(math.Round(long / 15.0))
	return time.FixedZone("Local", zoneOffset*3600)
}
