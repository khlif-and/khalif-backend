package domain

import (
	"time"

	"gorm.io/gorm"
)

// Surah represents a chapter in Al-Quran
type Surah struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SurahID     string         `gorm:"uniqueIndex;size:36;not null" json:"surah_id"` // UUID
	NameSurah   string         `gorm:"size:255;not null" json:"name_surah"`
	AudioURL    string         `gorm:"text" json:"audio_url"`
	NoSurah     int            `gorm:"not null" json:"no_surah"`
	TypeSurat   string         `gorm:"size:20" json:"type_surat"` // makiyah, madaniyah
	JuzNumber   int            `gorm:"not null" json:"juz_number"`
	TotalAyats  int            `gorm:"not null" json:"total_ayats"`
	TranslateEn string         `gorm:"type:text" json:"english_translate"`
	TranslateId string         `gorm:"type:text" json:"indonesian_translate"`
	
	Ayats       []Ayat         `gorm:"foreignKey:SurahID;references:ID" json:"ayats"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Ayat represents a verse in Al-Quran
type Ayat struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SurahID        uint           `gorm:"not null;index" json:"-"`
	AyatID         string         `gorm:"uniqueIndex;size:36;not null" json:"ayat_id"` // UUID
	NumberOfAyat   int            `gorm:"not null" json:"number_of_ayat"`
	AyatArabic     string         `gorm:"type:text" json:"ayat_arabic"`
	AyatLatin      string         `gorm:"type:text" json:"ayat_latin"`
	TranslateEn    string         `gorm:"type:text" json:"english_translate"`
	TranslateId    string         `gorm:"type:text" json:"indonesia_translate"`
	Tajwid         string         `gorm:"size:100" json:"tajwid"`       // idzar halqi, etc
	ColorTajwid    string         `gorm:"size:20" json:"color_tajwid"` // hex color code
	SajadahAyats   bool           `gorm:"default:false" json:"sajadah_ayats"`
	AudioURL       string         `gorm:"text" json:"audio_url"`
	AsbabunNuzul   string         `gorm:"type:text" json:"asbabun_nuzul"`

	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// DTOs for Request

// CreateAlquranRequest now uses standard JSON binding
type CreateAlquranRequest struct {
	SurahID       string `json:"surah_id"`
	NameSurah     string `json:"name_surah"`
	AudioURL      string `json:"audio_url"`
	
	// Translations for Surah
	EnglishTranslate    string `json:"english_translate"`
	IndonesianTranslate string `json:"indonesian_translate"`

	NoSurah     int    `json:"no_surah"`
	TypeSurat   string `json:"type_surat"`
	JuzNumber   int    `json:"juz_number"`
	TotalAyats  int    `json:"total_ayats"`

	// Parsed Ayats (Directly bound from JSON)
	Ayats       []AyatRequest `json:"ayats"`
}

type AyatRequest struct {
	AyatID         string       `json:"ayat_id"`
	NumberOfAyat   int          `json:"number_of_ayat"`
	AyatArabic     string       `json:"ayat_arabic"`
	AyatLatin      string       `json:"ayat_latin"`
	AyatTranslate  AyatTranslateRequest `json:"ayat_translate"`
	Tajwid         string       `json:"tajwid"`
	ColorTajwid    string       `json:"color_tajwid"`
	SajadahAyats   bool         `json:"sajadah_ayats"`
	AudioURL       string       `json:"audio_url"`
	AsbabunNuzul   string       `json:"asbabun_nuzul"`
}

type AyatTranslateRequest struct {
	Indonesia string `json:"indonesia"`
	English   string `json:"english"`
}

func (r *CreateAlquranRequest) ToEntity() Surah {
	// Construct Surah
	s := Surah{
		SurahID:     r.SurahID,
		NameSurah:   r.NameSurah,
		AudioURL:    r.AudioURL,
		NoSurah:     r.NoSurah,
		TypeSurat:   r.TypeSurat,
		JuzNumber:   r.JuzNumber,
		TotalAyats:  r.TotalAyats,
		TranslateEn: r.EnglishTranslate,
		TranslateId: r.IndonesianTranslate,
		Ayats:       make([]Ayat, 0, len(r.Ayats)),
	}

	for _, a := range r.Ayats {
		ayat := Ayat{
			AyatID:       a.AyatID,
			NumberOfAyat: a.NumberOfAyat,
			AyatArabic:   a.AyatArabic,
			AyatLatin:    a.AyatLatin,
			TranslateEn:  a.AyatTranslate.English,
			TranslateId:  a.AyatTranslate.Indonesia,
			Tajwid:       a.Tajwid,
			ColorTajwid:  a.ColorTajwid,
			SajadahAyats: a.SajadahAyats,
			AudioURL:     a.AudioURL,
			AsbabunNuzul: a.AsbabunNuzul,
		}
		s.Ayats = append(s.Ayats, ayat)
	}
	
	return s
}
