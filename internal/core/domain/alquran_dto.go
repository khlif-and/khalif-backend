package domain

// SurahListResponse for returning list of surahs
type SurahListResponse struct {
	Surahs []Surah `json:"surahs"`
	Total  int64   `json:"total"`
}
